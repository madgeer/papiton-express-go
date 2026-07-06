package main

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/payment/models"
)

// MockPaymentRepository implements IPaymentRepository for testing
type MockPaymentRepository struct {
	CreatePaymentFunc                 func(ctx context.Context, payment *models.Payment) error
	GetPaymentByIDFunc                func(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	UpdatePaymentStatusWithOutboxFunc func(ctx context.Context, payment *models.Payment, outbox *models.PaymentOutbox) error
}

func (m *MockPaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return m.CreatePaymentFunc(ctx, payment)
}

func (m *MockPaymentRepository) GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return m.GetPaymentByIDFunc(ctx, id)
}

func (m *MockPaymentRepository) UpdatePaymentStatusWithOutbox(ctx context.Context, payment *models.Payment, outbox *models.PaymentOutbox) error {
	return m.UpdatePaymentStatusWithOutboxFunc(ctx, payment, outbox)
}

// TestCreateInvoice_Success tests creating a payment invoice
func TestCreateInvoice_Success(t *testing.T) {
	mockRepo := &MockPaymentRepository{
		CreatePaymentFunc: func(ctx context.Context, payment *models.Payment) error {
			return nil
		},
	}

	service := NewPaymentService(mockRepo)

	req := &CreatePaymentRequest{
		OrderID:       uuid.New().String(),
		Amount:        45000.0,
		PaymentMethod: "GOPAY",
	}

	payment, err := service.CreateInvoice(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if payment.Status != models.StatusPending {
		t.Errorf("Expected status PENDING, got %s", payment.Status)
	}

	if payment.Amount != 45000.0 {
		t.Errorf("Expected amount 45000.0, got %f", payment.Amount)
	}
}

// TestProcessWebhook_Success verifies status update and outbox creation on payment success
func TestProcessWebhook_Success(t *testing.T) {
	paymentID := uuid.New()
	orderID := uuid.New()

	mockRepo := &MockPaymentRepository{
		GetPaymentByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
			return &models.Payment{
				ID:            paymentID,
				OrderID:       orderID,
				Amount:        25000.0,
				Status:        models.StatusPending,
				PaymentMethod: "OVO",
			}, nil
		},
		UpdatePaymentStatusWithOutboxFunc: func(ctx context.Context, payment *models.Payment, outbox *models.PaymentOutbox) error {
			if payment.Status != models.StatusSuccess {
				t.Errorf("Expected status SUCCESS, got %s", payment.Status)
			}
			if outbox.EventType != "PaymentVerified" {
				t.Errorf("Expected outbox event PaymentVerified, got %s", outbox.EventType)
			}
			return nil
		},
	}

	service := NewPaymentService(mockRepo)

	req := &WebhookPaymentRequest{
		PaymentID:     paymentID.String(),
		TransactionID: "TRX-12345",
		Status:        "SUCCESS",
	}

	payment, err := service.ProcessWebhook(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if payment.Status != models.StatusSuccess {
		t.Errorf("Expected payment status SUCCESS, got %s", payment.Status)
	}

	if *payment.TransactionID != "TRX-12345" {
		t.Errorf("Expected transaction ID TRX-12345, got %s", *payment.TransactionID)
	}
}

// TestProcessWebhook_AlreadyProcessed fails when trying to verify an already verified payment
func TestProcessWebhook_AlreadyProcessed(t *testing.T) {
	paymentID := uuid.New()

	mockRepo := &MockPaymentRepository{
		GetPaymentByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
			txID := "TRX-EXISTING"
			return &models.Payment{
				ID:            paymentID,
				Status:        models.StatusSuccess,
				TransactionID: &txID,
			}, nil
		},
	}

	service := NewPaymentService(mockRepo)

	req := &WebhookPaymentRequest{
		PaymentID:     paymentID.String(),
		TransactionID: "TRX-NEW",
		Status:        "SUCCESS",
	}

	_, err := service.ProcessWebhook(context.Background(), req)

	if err == nil {
		t.Fatalf("Expected error for already processed payment, got nil")
	}
}

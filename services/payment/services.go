package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/payment/models"
)

type PaymentService struct {
	repo IPaymentRepository
}

func NewPaymentService(repo IPaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

// CreateInvoice creates a new billing invoice in PENDING status
func (s *PaymentService) CreateInvoice(ctx context.Context, req *CreatePaymentRequest) (*models.Payment, error) {
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid orderId format")
	}

	payment := &models.Payment{
		ID:            uuid.New(),
		OrderID:       orderUUID,
		Amount:        req.Amount,
		Status:        models.StatusPending,
		PaymentMethod: req.PaymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPayment fetches payment details by ID
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return s.repo.GetPaymentByID(ctx, id)
}

// ProcessWebhook updates payment status and records an outbox event atomically
func (s *PaymentService) ProcessWebhook(ctx context.Context, req *WebhookPaymentRequest) (*models.Payment, error) {
	paymentUUID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid paymentId format")
	}

	payment, err := s.repo.GetPaymentByID(ctx, paymentUUID)
	if err != nil {
		return nil, fmt.Errorf("payment not found")
	}

	if payment.Status != models.StatusPending {
		return nil, fmt.Errorf("payment is already processed")
	}

	// Update fields
	payment.Status = models.PaymentStatus(req.Status)
	payment.TransactionID = &req.TransactionID
	payment.UpdatedAt = time.Now()

	// Prepare Event Type
	eventType := "PaymentFailed"
	if payment.Status == models.StatusSuccess {
		eventType = "PaymentVerified"
	}

	// Prepare Event Payload
	payloadBytes, err := json.Marshal(map[string]interface{}{
		"paymentId":     payment.ID.String(),
		"orderId":       payment.OrderID.String(),
		"amount":        payment.Amount,
		"transactionId": req.TransactionID,
		"status":        string(payment.Status),
	})
	if err != nil {
		return nil, errors.New("failed to serialize event payload")
	}

	outboxEvent := &models.PaymentOutbox{
		ID:            uuid.New(),
		AggregateType: "payment",
		AggregateID:   payment.ID,
		EventType:     eventType,
		Payload:       payloadBytes,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}

	// Save atomically inside a transaction
	err = s.repo.UpdatePaymentStatusWithOutbox(ctx, payment, outboxEvent)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

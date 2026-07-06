package main

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/models"
)

// MockOrderRepository implements IOrderRepository for testing purposes
type MockOrderRepository struct {
	GetTariffFunc          func(origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error)
	SaveOrderFunc          func(order *models.Order, addresses []models.Address, outbox *models.OrderOutbox) error
	GetOrderByIDFunc       func(id uuid.UUID) (*models.Order, error)
	ListOrdersFunc         func() ([]models.Order, error)
	UpdateOrderStatusFunc  func(id uuid.UUID, status models.OrderStatus) error
	GetOrderByIDSimpleFunc func(id uuid.UUID) (*models.Order, error)
}

func (m *MockOrderRepository) GetTariff(origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error) {
	return m.GetTariffFunc(origin, destination, serviceTypeID)
}

func (m *MockOrderRepository) SaveOrder(order *models.Order, addresses []models.Address, outbox *models.OrderOutbox) error {
	return m.SaveOrderFunc(order, addresses, outbox)
}

func (m *MockOrderRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	return m.GetOrderByIDFunc(id)
}

func (m *MockOrderRepository) ListOrders() ([]models.Order, error) {
	return m.ListOrdersFunc()
}

func (m *MockOrderRepository) UpdateOrderStatus(id uuid.UUID, status models.OrderStatus) error {
	return m.UpdateOrderStatusFunc(id, status)
}

func (m *MockOrderRepository) GetOrderByIDSimple(id uuid.UUID) (*models.Order, error) {
	return m.GetOrderByIDSimpleFunc(id)
}

// TestCreateOrder_Success tests the happy path of creating an order
func TestCreateOrder_Success(t *testing.T) {
	mockRepo := &MockOrderRepository{
		GetTariffFunc: func(origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error) {
			return &models.Tariff{
				PricePerKg: 10000.0,
			}, nil
		},
		SaveOrderFunc: func(order *models.Order, addresses []models.Address, outbox *models.OrderOutbox) error {
			return nil // Success
		},
	}

	service := NewOrderService(mockRepo)

	req := &CreateOrderRequest{
		CustomerID:    uuid.New().String(),
		ServiceTypeID: uuid.New().String(),
		Weight:        3.5,
		Length:        30,
		Width:         20,
		Height:        10,
	}
	req.Sender.Name = "Budi"
	req.Sender.City = "Jakarta"
	req.Receiver.Name = "Sari"
	req.Receiver.City = "Bandung"

	order, err := service.CreateOrder(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if order.Status != models.StatusWaitingPayment {
		t.Errorf("Expected status WAITING_PAYMENT, got %s", order.Status)
	}

	// weight (3.5) * tariff (10000) + insurance (5000) = 40000
	expectedTotalPrice := 40000.0
	if order.TotalPrice != expectedTotalPrice {
		t.Errorf("Expected total price %.2f, got %.2f", expectedTotalPrice, order.TotalPrice)
	}

	if order.TrackingNumber == "" {
		t.Errorf("Expected tracking number to be generated, got empty string")
	}
}

// TestCreateOrder_TariffNotFound tests when no tariff is found in the database
func TestCreateOrder_TariffNotFound(t *testing.T) {
	mockRepo := &MockOrderRepository{
		GetTariffFunc: func(origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error) {
			return nil, errors.New("record not found")
		},
	}

	service := NewOrderService(mockRepo)

	req := &CreateOrderRequest{
		CustomerID:    uuid.New().String(),
		ServiceTypeID: uuid.New().String(),
		Weight:        3.5,
	}
	req.Sender.City = "Jakarta"
	req.Receiver.City = "UnknownCity"

	_, err := service.CreateOrder(req)

	if err == nil {
		t.Fatalf("Expected error due to missing tariff, got nil")
	}
}

// TestCreateOrder_InvalidUUID tests validation of invalid UUID formats
func TestCreateOrder_InvalidUUID(t *testing.T) {
	mockRepo := &MockOrderRepository{}
	service := NewOrderService(mockRepo)

	req := &CreateOrderRequest{
		CustomerID:    "invalid-uuid",
		ServiceTypeID: uuid.New().String(),
		Weight:        3.5,
	}

	_, err := service.CreateOrder(req)

	if err == nil {
		t.Fatalf("Expected error due to invalid UUID, got nil")
	}
}

package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/dto"
	"github.com/madgeer/papiton-express-go/services/shipping/internal/service"
	"github.com/madgeer/papiton-express-go/services/shipping/models"
	"gorm.io/gorm"
)

type MockShippingRepository struct {
	CreateCourierFunc                   func(ctx context.Context, courier *models.Courier) error
	GetCourierByIDFunc                  func(ctx context.Context, id uuid.UUID) (*models.Courier, error)
	ListCouriersFunc                    func(ctx context.Context) ([]models.Courier, error)
	UpdateCourierStatusFunc             func(ctx context.Context, tx *gorm.DB, id uuid.UUID, status models.CourierStatus) error
	CreateShipmentWithCourierUpdateFunc func(ctx context.Context, shipment *models.Shipment) error
	GetShipmentByIDFunc                 func(ctx context.Context, id uuid.UUID) (*models.Shipment, error)
	UpdateShipmentStatusWithOutboxFunc  func(ctx context.Context, shipment *models.Shipment, outbox *models.ShippingOutbox, updateCourier bool, courierID uuid.UUID, courierStatus models.CourierStatus) error
}

func (m *MockShippingRepository) CreateCourier(ctx context.Context, courier *models.Courier) error {
	return m.CreateCourierFunc(ctx, courier)
}

func (m *MockShippingRepository) GetCourierByID(ctx context.Context, id uuid.UUID) (*models.Courier, error) {
	return m.GetCourierByIDFunc(ctx, id)
}

func (m *MockShippingRepository) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	return m.ListCouriersFunc(ctx)
}

func (m *MockShippingRepository) UpdateCourierStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status models.CourierStatus) error {
	return m.UpdateCourierStatusFunc(ctx, tx, id, status)
}

func (m *MockShippingRepository) CreateShipmentWithCourierUpdate(ctx context.Context, shipment *models.Shipment) error {
	return m.CreateShipmentWithCourierUpdateFunc(ctx, shipment)
}

func (m *MockShippingRepository) GetShipmentByID(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	return m.GetShipmentByIDFunc(ctx, id)
}

func (m *MockShippingRepository) UpdateShipmentStatusWithOutbox(
	ctx context.Context,
	shipment *models.Shipment,
	outbox *models.ShippingOutbox,
	updateCourier bool,
	courierID uuid.UUID,
	courierStatus models.CourierStatus,
) error {
	return m.UpdateShipmentStatusWithOutboxFunc(ctx, shipment, outbox, updateCourier, courierID, courierStatus)
}

func TestCreateCourier_Success(t *testing.T) {
	mockRepo := &MockShippingRepository{
		CreateCourierFunc: func(ctx context.Context, courier *models.Courier) error {
			return nil
		},
	}

	svc := service.NewShippingService(mockRepo)

	req := &dto.CreateCourierRequest{
		Name:  "Anto",
		Phone: "085544332211",
	}

	courier, err := svc.CreateCourier(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if courier.Status != models.CourierAvailable {
		t.Errorf("Expected status AVAILABLE, got %s", courier.Status)
	}
}

func TestAssignCourier_Success(t *testing.T) {
	courierID := uuid.New()
	orderID := uuid.New()

	mockRepo := &MockShippingRepository{
		GetCourierByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Courier, error) {
			return &models.Courier{
				ID:     courierID,
				Name:   "Anto",
				Status: models.CourierAvailable,
			}, nil
		},
		CreateShipmentWithCourierUpdateFunc: func(ctx context.Context, shipment *models.Shipment) error {
			if shipment.OrderID != orderID {
				t.Errorf("Expected orderID %s, got %s", orderID, shipment.OrderID)
			}
			return nil
		},
	}

	svc := service.NewShippingService(mockRepo)

	req := &dto.AssignShipmentRequest{
		OrderID:   orderID.String(),
		CourierID: courierID.String(),
	}

	shipment, err := svc.AssignCourier(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shipment.Status != models.ShipmentPendingPickup {
		t.Errorf("Expected shipment status PENDING_PICKUP, got %s", shipment.Status)
	}
}

func TestAssignCourier_CourierBusy(t *testing.T) {
	courierID := uuid.New()

	mockRepo := &MockShippingRepository{
		GetCourierByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Courier, error) {
			return &models.Courier{
				ID:     courierID,
				Name:   "Anto",
				Status: models.CourierOnDelivery, // Courier busy!
			}, nil
		},
	}

	svc := service.NewShippingService(mockRepo)

	req := &dto.AssignShipmentRequest{
		OrderID:   uuid.New().String(),
		CourierID: courierID.String(),
	}

	_, err := svc.AssignCourier(context.Background(), req)

	if err == nil {
		t.Fatalf("Expected error due to busy courier, got nil")
	}
}

func TestUpdateShipmentStatus_Delivered(t *testing.T) {
	shipmentID := uuid.New()
	courierID := uuid.New()

	mockRepo := &MockShippingRepository{
		GetShipmentByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
			return &models.Shipment{
				ID:        shipmentID,
				OrderID:   uuid.New(),
				CourierID: &courierID,
				Status:    models.ShipmentPickedUp,
			}, nil
		},
		UpdateShipmentStatusWithOutboxFunc: func(ctx context.Context, shipment *models.Shipment, outbox *models.ShippingOutbox, updateCourier bool, cID uuid.UUID, cStatus models.CourierStatus) error {
			if !updateCourier {
				t.Errorf("Expected updateCourier to be true")
			}
			if cID != courierID {
				t.Errorf("Expected courierID %s, got %s", courierID, cID)
			}
			if cStatus != models.CourierAvailable {
				t.Errorf("Expected courier status AVAILABLE, got %s", cStatus)
			}
			if outbox.EventType != "ShipmentStatusUpdated" {
				t.Errorf("Expected event Type ShipmentStatusUpdated, got %s", outbox.EventType)
			}
			return nil
		},
	}

	svc := service.NewShippingService(mockRepo)

	req := &dto.UpdateShipmentStatusRequest{
		Status: "DELIVERED",
		Notes:  "Package delivered to receiver Budi",
	}

	shipment, err := svc.UpdateShipmentStatus(context.Background(), shipmentID, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shipment.Status != models.ShipmentDelivered {
		t.Errorf("Expected status DELIVERED, got %s", shipment.Status)
	}
}

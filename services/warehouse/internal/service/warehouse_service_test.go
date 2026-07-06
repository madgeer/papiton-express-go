package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/dto"
	"github.com/madgeer/papiton-express-go/services/warehouse/internal/service"
	"github.com/madgeer/papiton-express-go/services/warehouse/models"
)

type MockWarehouseRepository struct {
	CreateWarehouseFunc       func(ctx context.Context, wh *models.Warehouse) error
	GetWarehouseByIDFunc      func(ctx context.Context, id uuid.UUID) (*models.Warehouse, error)
	ListWarehousesFunc        func(ctx context.Context) ([]models.Warehouse, error)
	CreateRouteFunc           func(ctx context.Context, route *models.WarehouseRoute) error
	GetNextWarehouseRouteFunc func(ctx context.Context, originID uuid.UUID, destinationCity string) (*models.WarehouseRoute, error)
	SaveMovementWithOutboxFunc func(ctx context.Context, movement *models.WarehouseMovement, outbox *models.WarehouseOutbox) error
}

func (m *MockWarehouseRepository) CreateWarehouse(ctx context.Context, wh *models.Warehouse) error {
	return m.CreateWarehouseFunc(ctx, wh)
}

func (m *MockWarehouseRepository) GetWarehouseByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error) {
	return m.GetWarehouseByIDFunc(ctx, id)
}

func (m *MockWarehouseRepository) ListWarehouses(ctx context.Context) ([]models.Warehouse, error) {
	return m.ListWarehousesFunc(ctx)
}

func (m *MockWarehouseRepository) CreateRoute(ctx context.Context, route *models.WarehouseRoute) error {
	return m.CreateRouteFunc(ctx, route)
}

func (m *MockWarehouseRepository) GetNextWarehouseRoute(ctx context.Context, originID uuid.UUID, destinationCity string) (*models.WarehouseRoute, error) {
	return m.GetNextWarehouseRouteFunc(ctx, originID, destinationCity)
}

func (m *MockWarehouseRepository) SaveMovementWithOutbox(ctx context.Context, movement *models.WarehouseMovement, outbox *models.WarehouseOutbox) error {
	return m.SaveMovementWithOutboxFunc(ctx, movement, outbox)
}

func TestCreateWarehouse_Success(t *testing.T) {
	mockRepo := &MockWarehouseRepository{
		CreateWarehouseFunc: func(ctx context.Context, wh *models.Warehouse) error {
			return nil
		},
	}

	svc := service.NewWarehouseService(mockRepo)

	req := &dto.CreateWarehouseRequest{
		Name: "Gudang Jakarta",
		City: "Jakarta",
	}

	wh, err := svc.CreateWarehouse(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if wh.Name != "Gudang Jakarta" {
		t.Errorf("Expected name 'Gudang Jakarta', got %s", wh.Name)
	}
}

func TestRecordMovementIn_WithRoutingRule(t *testing.T) {
	whID := uuid.New()
	nextWhID := uuid.New()
	orderID := uuid.New()

	mockRepo := &MockWarehouseRepository{
		GetWarehouseByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Warehouse, error) {
			return &models.Warehouse{
				ID:   whID,
				Name: "Gudang Bandung",
				City: "Bandung",
			}, nil
		},
		GetNextWarehouseRouteFunc: func(ctx context.Context, originID uuid.UUID, destCity string) (*models.WarehouseRoute, error) {
			return &models.WarehouseRoute{
				OriginWarehouseID: whID,
				DestinationCity:   "Surabaya",
				NextWarehouseID:   &nextWhID,
			}, nil
		},
		SaveMovementWithOutboxFunc: func(ctx context.Context, m *models.WarehouseMovement, o *models.WarehouseOutbox) error {
			if m.Type != models.MovementIn {
				t.Errorf("Expected type WAREHOUSE_IN, got %s", m.Type)
			}
			if o.EventType != "PackageWarehouseIn" {
				t.Errorf("Expected event PackageWarehouseIn, got %s", o.EventType)
			}
			return nil
		},
	}

	svc := service.NewWarehouseService(mockRepo)

	req := &dto.CreateMovementRequest{
		OrderID:         orderID.String(),
		WarehouseID:     whID.String(),
		Type:            "WAREHOUSE_IN",
		DestinationCity: "Surabaya",
		Notes:           "Transit package",
	}

	res, err := svc.RecordMovement(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.NextWarehouseID == nil || *res.NextWarehouseID != nextWhID.String() {
		t.Errorf("Expected next warehouse %s, got %v", nextWhID, res.NextWarehouseID)
	}
}

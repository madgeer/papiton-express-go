package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/warehouse/models"
	"gorm.io/gorm"
)

type IWarehouseRepository interface {
	CreateWarehouse(ctx context.Context, wh *models.Warehouse) error
	GetWarehouseByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error)
	ListWarehouses(ctx context.Context) ([]models.Warehouse, error)
	CreateRoute(ctx context.Context, route *models.WarehouseRoute) error
	GetNextWarehouseRoute(ctx context.Context, originID uuid.UUID, destinationCity string) (*models.WarehouseRoute, error)
	SaveMovementWithOutbox(ctx context.Context, movement *models.WarehouseMovement, outbox *models.WarehouseOutbox) error
}

type WarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) IWarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) CreateWarehouse(ctx context.Context, wh *models.Warehouse) error {
	return r.db.WithContext(ctx).Create(wh).Error
}

func (r *WarehouseRepository) GetWarehouseByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error) {
	var wh models.Warehouse
	err := r.db.WithContext(ctx).First(&wh, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func (r *WarehouseRepository) ListWarehouses(ctx context.Context) ([]models.Warehouse, error) {
	var list []models.Warehouse
	err := r.db.WithContext(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *WarehouseRepository) CreateRoute(ctx context.Context, route *models.WarehouseRoute) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *WarehouseRepository) GetNextWarehouseRoute(ctx context.Context, originID uuid.UUID, destinationCity string) (*models.WarehouseRoute, error) {
	var route models.WarehouseRoute
	err := r.db.WithContext(ctx).Where(
		"origin_warehouse_id = ? AND LOWER(destination_city) = LOWER(?)",
		originID,
		destinationCity,
	).First(&route).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *WarehouseRepository) SaveMovementWithOutbox(ctx context.Context, movement *models.WarehouseMovement, outbox *models.WarehouseOutbox) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(movement).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}

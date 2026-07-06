package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/shipping/models"
	"gorm.io/gorm"
)

type IShippingRepository interface {
	CreateCourier(ctx context.Context, courier *models.Courier) error
	GetCourierByID(ctx context.Context, id uuid.UUID) (*models.Courier, error)
	ListCouriers(ctx context.Context) ([]models.Courier, error)
	UpdateCourierStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status models.CourierStatus) error
	CreateShipmentWithCourierUpdate(ctx context.Context, shipment *models.Shipment) error
	GetShipmentByID(ctx context.Context, id uuid.UUID) (*models.Shipment, error)
	UpdateShipmentStatusWithOutbox(ctx context.Context, shipment *models.Shipment, outbox *models.ShippingOutbox, updateCourier bool, courierID uuid.UUID, courierStatus models.CourierStatus) error
}

type ShippingRepository struct {
	db *gorm.DB
}

func NewShippingRepository(db *gorm.DB) IShippingRepository {
	return &ShippingRepository{db: db}
}

func (r *ShippingRepository) CreateCourier(ctx context.Context, courier *models.Courier) error {
	return r.db.WithContext(ctx).Create(courier).Error
}

func (r *ShippingRepository) GetCourierByID(ctx context.Context, id uuid.UUID) (*models.Courier, error) {
	var courier models.Courier
	err := r.db.WithContext(ctx).First(&courier, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &courier, nil
}

func (r *ShippingRepository) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	var couriers []models.Courier
	err := r.db.WithContext(ctx).Find(&couriers).Error
	if err != nil {
		return nil, err
	}
	return couriers, nil
}

func (r *ShippingRepository) UpdateCourierStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status models.CourierStatus) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Model(&models.Courier{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ShippingRepository) CreateShipmentWithCourierUpdate(ctx context.Context, shipment *models.Shipment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(shipment).Error; err != nil {
			return err
		}
		if shipment.CourierID != nil {
			if err := r.UpdateCourierStatus(ctx, tx, *shipment.CourierID, models.CourierOnDelivery); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ShippingRepository) GetShipmentByID(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.WithContext(ctx).Preload("Courier").First(&shipment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *ShippingRepository) UpdateShipmentStatusWithOutbox(
	ctx context.Context,
	shipment *models.Shipment,
	outbox *models.ShippingOutbox,
	updateCourier bool,
	courierID uuid.UUID,
	courierStatus models.CourierStatus,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(shipment).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		if updateCourier {
			if err := r.UpdateCourierStatus(ctx, tx, courierID, courierStatus); err != nil {
				return err
			}
		}
		return nil
	})
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/models"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	GetTariff(ctx context.Context, origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error)
	SaveOrder(ctx context.Context, order *models.Order, addresses []models.Address, outbox *models.OrderOutbox) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	ListOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error
	GetOrderByIDSimple(ctx context.Context, id uuid.UUID) (*models.Order, error)
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetTariff(ctx context.Context, origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error) {
	var tariff models.Tariff
	err := r.db.WithContext(ctx).Where(
		"LOWER(origin_city) = LOWER(?) AND LOWER(destination_city) = LOWER(?) AND service_type_id = ?",
		origin,
		destination,
		serviceTypeID,
	).First(&tariff).Error
	if err != nil {
		return nil, err
	}
	return &tariff, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order *models.Order, addresses []models.Address, outbox *models.OrderOutbox) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, addr := range addresses {
			if err := tx.Create(&addr).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Preload("Addresses").Preload("ServiceType").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).Preload("Addresses").Preload("ServiceType").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) GetOrderByIDSimple(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

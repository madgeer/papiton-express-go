package main

import (
	"github.com/google/uuid"
	"github.com/madgeer/papiton-express-go/services/order/models"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	GetTariff(origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error)
	SaveOrder(order *models.Order, addresses []models.Address, outbox *models.OrderOutbox) error
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	ListOrders() ([]models.Order, error)
	UpdateOrderStatus(id uuid.UUID, status models.OrderStatus) error
	GetOrderByIDSimple(id uuid.UUID) (*models.Order, error)
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

// GetTariff finds the shipping price per kg for a specific city route and service type
func (r *OrderRepository) GetTariff(origin, destination string, serviceTypeID uuid.UUID) (*models.Tariff, error) {
	var tariff models.Tariff
	err := r.db.Where(
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

// SaveOrder saves the order, sender/receiver addresses, and outbox event inside a single transaction
func (r *OrderRepository) SaveOrder(order *models.Order, addresses []models.Address, outbox *models.OrderOutbox) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

// GetOrderByID retrieves a single order preloaded with addresses and service type details
func (r *OrderRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Addresses").Preload("ServiceType").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListOrders retrieves all orders in the system
func (r *OrderRepository) ListOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Addresses").Preload("ServiceType").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateOrderStatus updates only the status column of an order
func (r *OrderRepository) UpdateOrderStatus(id uuid.UUID, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

// GetOrderByIDSimple retrieves an order without any relationship preloading
func (r *OrderRepository) GetOrderByIDSimple(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

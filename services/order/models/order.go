package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusWaitingPayment OrderStatus = "WAITING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusCourierAssign  OrderStatus = "COURIER_ASSIGNED"
	StatusPickedUp       OrderStatus = "PICKED_UP"
	StatusWarehouseIn    OrderStatus = "WAREHOUSE_IN"
	StatusInDelivery     OrderStatus = "IN_DELIVERY"
	StatusDelivered      OrderStatus = "DELIVERED"
	StatusCompleted      OrderStatus = "COMPLETED"
)

type Order struct {
	ID             uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID     uuid.UUID   `gorm:"type:uuid;not null" json:"customerId"`
	ServiceTypeID  uuid.UUID   `gorm:"type:uuid;not null" json:"serviceTypeId"`
	ServiceType    ServiceType `gorm:"foreignKey:ServiceTypeID" json:"serviceType,omitempty"`
	TrackingNumber string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"trackingNumber"`
	Weight         float64     `gorm:"type:numeric(10,2);not null" json:"weight"`
	Length         float64     `gorm:"type:numeric(10,2);not null" json:"length"`
	Width          float64     `gorm:"type:numeric(10,2);not null" json:"width"`
	Height         float64     `gorm:"type:numeric(10,2);not null" json:"height"`
	ShippingCost   float64     `gorm:"type:numeric(12,2);not null" json:"shippingCost"`
	InsuranceFee   float64     `gorm:"type:numeric(12,2);not null" json:"insuranceFee"`
	TotalPrice     float64     `gorm:"type:numeric(12,2);not null" json:"totalPrice"`
	Status         OrderStatus `gorm:"type:order_status;default:WAITING_PAYMENT;not null" json:"status"`
	Addresses      []Address   `gorm:"foreignKey:OrderID" json:"addresses,omitempty"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

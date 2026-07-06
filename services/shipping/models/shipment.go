package models

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus string

const (
	ShipmentPendingPickup  ShipmentStatus = "PENDING_PICKUP"
	ShipmentPickedUp       ShipmentStatus = "PICKED_UP"
	ShipmentOnTransit      ShipmentStatus = "ON_TRANSIT"
	ShipmentOutForDelivery ShipmentStatus = "OUT_FOR_DELIVERY"
	ShipmentDelivered      ShipmentStatus = "DELIVERED"
	ShipmentFailed         ShipmentStatus = "FAILED"
)

type Shipment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID   uuid.UUID      `gorm:"type:uuid;not null" json:"orderId"`
	CourierID *uuid.UUID     `gorm:"type:uuid" json:"courierId,omitempty"`
	Status    ShipmentStatus `gorm:"type:shipment_status;default:'PENDING_PICKUP';not null" json:"status"`
	Notes     string         `gorm:"type:text" json:"notes,omitempty"`
	Courier   *Courier       `gorm:"foreignKey:CourierID" json:"courier,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

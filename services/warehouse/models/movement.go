package models

import (
	"time"

	"github.com/google/uuid"
)

type MovementType string

const (
	MovementIn  MovementType = "WAREHOUSE_IN"
	MovementOut MovementType = "WAREHOUSE_OUT"
)

type WarehouseMovement struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID     uuid.UUID    `gorm:"type:uuid;not null" json:"orderId"`
	WarehouseID uuid.UUID    `gorm:"type:uuid;not null" json:"warehouseId"`
	Type        MovementType `gorm:"type:movement_type;not null" json:"type"`
	Notes       string       `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
}

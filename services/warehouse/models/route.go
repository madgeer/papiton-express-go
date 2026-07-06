package models

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseRoute struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	OriginWarehouseID uuid.UUID  `gorm:"type:uuid;not null" json:"originWarehouseId"`
	DestinationCity   string     `gorm:"type:varchar(100);not null" json:"destinationCity"`
	NextWarehouseID   *uuid.UUID `gorm:"type:uuid" json:"nextWarehouseId,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

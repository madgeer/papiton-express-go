package models

import (
	"time"

	"github.com/google/uuid"
)

type Tariff struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	OriginCity      string      `gorm:"type:varchar(100);not null" json:"originCity"`
	DestinationCity string      `gorm:"type:varchar(100);not null" json:"destinationCity"`
	ServiceTypeID   uuid.UUID   `gorm:"type:uuid;not null" json:"serviceTypeId"`
	ServiceType     ServiceType `gorm:"foreignKey:ServiceTypeID" json:"serviceType,omitempty"`
	PricePerKg      float64     `gorm:"type:numeric(12,2);not null" json:"pricePerKg"`
	CreatedAt       time.Time   `json:"createdAt"`
}

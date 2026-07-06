package models

import (
	"time"

	"github.com/google/uuid"
)

type CourierStatus string

const (
	CourierAvailable  CourierStatus = "AVAILABLE"
	CourierOnDelivery CourierStatus = "ON_DELIVERY"
	CourierOffline    CourierStatus = "OFFLINE"
)

type Courier struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string        `gorm:"type:varchar(100);not null" json:"name"`
	Phone     string        `gorm:"type:varchar(20);not null" json:"phone"`
	Status    CourierStatus `gorm:"type:courier_status;default:'AVAILABLE';not null" json:"status"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
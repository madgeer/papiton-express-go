package models

import (
	"time"

	"github.com/google/uuid"
)

type Tracking struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID        uuid.UUID       `gorm:"type:uuid;uniqueIndex;not null" json:"orderId"`
	TrackingNumber string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"trackingNumber"`
	CurrentStatus  string          `gorm:"type:varchar(50);not null" json:"currentStatus"`
	LastUpdated    time.Time       `json:"lastUpdated"`
	Events         []TrackingEvent `gorm:"foreignKey:TrackingID;constraint:OnDelete:CASCADE" json:"events,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

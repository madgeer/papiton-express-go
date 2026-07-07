package models

import (
	"time"

	"github.com/google/uuid"
)

type TrackingEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TrackingID  uuid.UUID `gorm:"type:uuid;not null" json:"trackingId"`
	Status      string    `gorm:"type:varchar(50);not null" json:"status"`
	Location    string    `gorm:"type:varchar(100)" json:"location,omitempty"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

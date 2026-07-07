package models

import (
	"time"

	"github.com/google/uuid"
)

type TrackingOutbox struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	AggregateType string    `gorm:"type:varchar(50);not null"`
	AggregateID   uuid.UUID `gorm:"type:uuid;not null"`
	EventType     string    `gorm:"type:varchar(100);not null"`
	Payload       []byte    `gorm:"type:jsonb;not null"`
	Status        string    `gorm:"type:varchar(20);default:'PENDING';not null"`
	CreatedAt     time.Time
	ProcessedAt   *time.Time
}

func (TrackingOutbox) TableName() string {
	return "tracking_outbox"
}

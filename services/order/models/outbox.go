package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OrderOutbox struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	AggregateType string          `gorm:"type:varchar(100);not null" json:"aggregateType"`
	AggregateID   uuid.UUID       `gorm:"type:uuid;not null" json:"aggregateId"`
	EventType     string          `gorm:"type:varchar(100);not null" json:"eventType"`
	Payload       json.RawMessage `gorm:"type:jsonb;not null" json:"payload"`
	Status        string          `gorm:"type:varchar(50);default:PENDING;not null" json:"status"`
	CreatedAt     time.Time       `json:"createdAt"`
	ProcessedAt   *time.Time      `gorm:"default:null" json:"processedAt,omitempty"`
}

func (OrderOutbox) TableName() string {
	return "order_outbox"
}

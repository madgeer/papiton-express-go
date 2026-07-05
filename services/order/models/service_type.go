package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceType struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	EstimatedDay int       `gorm:"not null" json:"estimatedDay"`
	CreatedAt    time.Time `json:"createdAt"`
}

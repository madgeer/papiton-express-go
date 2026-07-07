package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null" json:"orderId"`
	RecipientEmail string    `gorm:"type:varchar(100);not null" json:"recipientEmail"`
	RecipientPhone string    `gorm:"type:varchar(20);not null" json:"recipientPhone"`
	Type           string    `gorm:"type:varchar(50);not null" json:"type"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	Status         string    `gorm:"type:varchar(20);default:'SENT';not null" json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

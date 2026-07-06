package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusPending PaymentStatus = "PENDING"
	StatusSuccess PaymentStatus = "SUCCESS"
	StatusFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID       uuid.UUID     `gorm:"type:uuid;not null" json:"orderId"`
	Amount        float64       `gorm:"type:decimal(15,2);not null" json:"amount"`
	Status        PaymentStatus `gorm:"type:payment_status;default:'PENDING';not null" json:"status"`
	PaymentMethod string        `gorm:"type:varchar(50);not null" json:"paymentMethod"`
	TransactionID *string       `gorm:"type:varchar(100)" json:"transactionId,omitempty"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

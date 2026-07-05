package models

import (
	"time"

	"github.com/google/uuid"
)

type AddressType string

const (
	AddressSender   AddressType = "SENDER"
	AddressReceiver AddressType = "RECEIVER"
)

type Address struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID    uuid.UUID   `gorm:"type:uuid;not null" json:"orderId"`
	Type       AddressType `gorm:"type:address_type;not null" json:"type"`
	Name       string      `gorm:"type:varchar(100);not null" json:"name"`
	Phone      string      `gorm:"type:varchar(20);not null" json:"phone"`
	Address    string      `gorm:"type:text;not null" json:"address"`
	City       string      `gorm:"type:varchar(100);not null" json:"city"`
	Province   string      `gorm:"type:varchar(100);not null" json:"province"`
	PostalCode string      `gorm:"type:varchar(20);not null" json:"postalCode"`
	CreatedAt  time.Time   `json:"createdAt"`
}

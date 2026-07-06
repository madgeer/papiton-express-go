package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin          UserRole = "ADMIN"
	RoleCustomer       UserRole = "CUSTOMER"
	RoleCourier        UserRole = "COURIER"
	RoleWarehouseStaff UserRole = "WAREHOUSE_STAFF"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Phone     string    `gorm:"type:varchar(20);not null" json:"phone"`
	Role      UserRole  `gorm:"type:user_role;default:'CUSTOMER';not null" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

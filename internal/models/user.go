package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	BaseModel
	Username string   `gorm:"not null;uniqueIndex" json:"username" validate:"required,min=2,max=30"`
	Email    string   `gorm:"not null;unique" json:"email" validate:"required,email"`
	Password string   `gorm:"not null" json:"password" validate:"required,min=8"`
	Role     UserRole `gorm:"type:varchar(50);default:'user'" json:"role"`
}

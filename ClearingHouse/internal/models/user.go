package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID        string `gorm:"type:char(36);primary_key;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Username  string     `gorm:"size:255;not null;uniqueIndex" json:"username" validate:"required,min=2,max=30"`
	Email     string     `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string     `gorm:"not null" json:"password" validate:"required,min=8"`
	Role      UserRole   `gorm:"type:varchar(50);default:'user'" json:"role"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	return
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if !u.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	var existingUser User
	if err := tx.Unscoped().Where("id = ?", u.ID).First(&existingUser).Error; err != nil {
		return err
	}

	if u.Role != existingUser.Role {
		return errors.New("role cannot be changed")
	}

	return nil
}
func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleUser, UserRoleAdmin:
		return true
	}
	return false
}

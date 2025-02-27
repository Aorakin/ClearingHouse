package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}

type User struct {
	BaseModel
	Username        string       `gorm:"not null;uniqueIndex" json:"username" validate:"required,min=2,max=30"`
	Email           string       `gorm:"not null;unique" json:"email" validate:"required,email"`
	Password        string       `gorm:"not null" json:"password" validate:"required,min=8"`
	Role            UserRole     `gorm:"type:varchar(50);default:'user'" json:"role"`
	Tickets         []Ticket     `gorm:"foreignKey:RequesterID"  json:"-"`
	OwnedProjects   []*Project   `gorm:"many2many:project_owners;" json:"owned_projects"`
	MemberProjects  []*Project   `gorm:"many2many:project_members;" json:"member_projects"`
	MemberNamespace []*Namespace `gorm:"many2many:namespace_members;" json:"member_namespaces"`
}

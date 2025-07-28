package models

import "github.com/google/uuid"

type Namespace struct {
	BaseModel
	Name        string    `gorm:"not null;uniqueIndex" json:"name" validate:"required,min=2,max=100"`
	Description string    `gorm:"not null" json:"description" validate:"required,min=2,max=500"`
	ProjectID   uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	Project     Project   `gorm:"foreignKey:ProjectID" json:"project"`
	Members     []User    `gorm:"many2many:namespace_members;" json:"members"`
}

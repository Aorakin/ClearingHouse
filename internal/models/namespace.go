package models

import "github.com/google/uuid"

type Namespace struct {
	BaseModel
	Name        string           `json:"name" `
	Description string           `json:"description"`
	Credit      float32          `gorm:"not null" json:"credit" validate:"required"`
	ProjectID   uuid.UUID        `gorm:"type:uuid;not null" json:"project_id"`
	Project     Project          `gorm:"foreignKey:ProjectID" json:"-"`
	Quotas      []NamespaceQuota `gorm:"many2many:namespace_quotas;" json:"quotas"`
	Members     []User           `gorm:"many2many:namespace_members;" json:"namespace_members"`
	OwnerID     uuid.UUID        `gorm:"type:uuid" json:"owner_id"`
}

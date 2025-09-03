package models

import "github.com/google/uuid"

type Namespace struct {
	BaseModel
	Name         string              `gorm:"not null" json:"name" validate:"required,min=2,max=100"`
	Description  string              `gorm:"not null" json:"description" validate:"required,min=2,max=500"`
	Credit       float32             `gorm:"not null" json:"credit" validate:"required"`
	ProjectID    uuid.UUID           `gorm:"type:uuid;not null" json:"project_id"`
	Project      Project             `gorm:"foreignKey:ProjectID" json:"-"`
	QuotaGroupID *uuid.UUID          `gorm:"type:uuid" json:"quota_group_id"`
	QuotaGroup   NamespaceQuotaGroup `gorm:"foreignKey:QuotaGroupID" json:"-"`
	Members      []User              `gorm:"many2many:namespace_members;" json:"members"`
}

// type Namespace struct {
// 	BaseModel
// 	Name        string           `json:"name" `
// 	Description string           `json:"description"`
// 	Credit      float32          `gorm:"not null" json:"credit" validate:"required"`
// 	ProjectID   uuid.UUID        `gorm:"type:uuid;not null" json:"project_id"`
// 	Project     Project          `gorm:"foreignKey:ProjectID" json:"-"`
// 	Quotas      []NamespaceQuota `gorm:"many2many:namespace_quotas;" json:"quotas"`
// 	Members     []User           `gorm:"many2many:project_members;" json:"project_members"`
// }

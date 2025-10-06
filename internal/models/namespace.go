package models

import "github.com/google/uuid"

type Namespace struct {
	BaseModel
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Credit      float32          `json:"credit"`
	ProjectID   *uuid.UUID       `gorm:"type:uuid" json:"project_id"`
	Project     *Project         `gorm:"foreignKey:ProjectID" json:"-"`
	OrgID       uuid.UUID        `gorm:"type:uuid" json:"organization_id"`
	Quotas      []NamespaceQuota `gorm:"many2many:namespace_quotas;" json:"quotas"`
	Members     []User           `gorm:"many2many:namespace_members;" json:"namespace_members"`
	Owner       *User            `gorm:"foreignKey:NamespaceID;references:ID" json:"owner"`
}

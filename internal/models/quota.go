package models

import "github.com/google/uuid"

type ResourceProperty struct {
	BaseModel
	ResourceID  uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	Resource    Resource  `gorm:"foreignKey:ResourceID" json:"-"`
	Price       float32   `gorm:"not null" json:"price"`
	MaxDuration float32   `json:"max_duration"`
}

type ResourceQuantity struct {
	BaseModel
	OrganizationQuotaID *uuid.UUID       `gorm:"type:uuid;index" json:"organization_quota_id,omitempty"`
	ProjectQuotaID      *uuid.UUID       `gorm:"type:uuid;index" json:"project_quota_id,omitempty"`
	NamespaceQuotaID    *uuid.UUID       `gorm:"type:uuid;index" json:"namespace_quota_id,omitempty"`
	Quantity            uint             `json:"quantity"`
	ResourcePropID      uuid.UUID        `gorm:"type:uuid;not null" json:"resource_property_id"`
	ResourceProp        ResourceProperty `gorm:"foreignKey:ResourcePropID" json:"resource_prop"`
}

type OrganizationQuota struct {
	BaseModel
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	ResourcePoolID uuid.UUID          `gorm:"type:uuid;not null" json:"resource_pool_id"`
	FromOrgID      uuid.UUID          `gorm:"type:uuid;not null" json:"from_organization_id"`
	ToOrgID        uuid.UUID          `gorm:"type:uuid;not null" json:"to_organization_id"`
	FromOrg        Organization       `gorm:"foreignKey:FromOrgID" json:"-"`
	ToOrg          Organization       `gorm:"foreignKey:ToOrgID" json:"-"`
	Resources      []ResourceQuantity `gorm:"foreignKey:OrganizationQuotaID" json:"resources"`
}

type ProjectQuota struct {
	BaseModel
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	OrganizationID      uuid.UUID          `gorm:"type:uuid" json:"organization_id"`
	OrganizationQuotaID uuid.UUID          `gorm:"type:uuid;not null" json:"organization_quota_id"`
	ProjectID           uuid.UUID          `gorm:"type:uuid;not null" json:"project_id"`
	ResourcePoolID      uuid.UUID          `gorm:"type:uuid;not null" json:"resource_pool_id"`
	Project             Project            `gorm:"foreignKey:ProjectID" json:"-"`
	Resources           []ResourceQuantity `gorm:"foreignKey:ProjectQuotaID" json:"resources"`
}

type NamespaceQuota struct {
	BaseModel
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	ProjectID      uuid.UUID          `gorm:"type:uuid" json:"project_id"`
	ProjectQuotaID uuid.UUID          `gorm:"type:uuid;not null" json:"project_quota_id"`
	ResourcePoolID uuid.UUID          `gorm:"type:uuid;not null" json:"resource_pool_id"`
	Namespaces     []Namespace        `gorm:"many2many:namespace_quotas;" json:"-"`
	Resources      []ResourceQuantity `gorm:"foreignKey:NamespaceQuotaID" json:"resources"`
}

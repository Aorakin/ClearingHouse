package models

import "github.com/google/uuid"

type OrganizationQuotaGroup struct {
	BaseModel
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	Resources          []ResourceQuantity `gorm:"foreignKey:OrganizationQuotaGroupID" json:"resources"`
	FromOrganizationID uuid.UUID          `gorm:"type:uuid;not null" json:"from_organization_id"`
	ToOrganizationID   uuid.UUID          `gorm:"type:uuid;not null" json:"to_organization_id"`
	FromOrganization   Organization       `gorm:"foreignKey:FromOrganizationID" json:"from_organization"`
	ToOrganization     Organization       `gorm:"foreignKey:ToOrganizationID" json:"to_organization"`
}

type ProjectQuotaGroup struct {
	BaseModel
	Name                     string             `json:"name"`
	Description              string             `json:"description"`
	OrganizationID           uuid.UUID          `gorm:"type:uuid;not null" json:"organization_id"`
	OrganizationQuotaGroupID uuid.UUID          `gorm:"type:uuid;not null" json:"organization_quota_group_id"`
	Resources                []ResourceQuantity `gorm:"foreignKey:ProjectQuotaGroupID" json:"resources"`
	Projects                 []Project          `gorm:"many2many:project_quotas;" json:"projects"`
	Organization             Organization       `gorm:"foreignKey:OrganizationID" json:"organization"`
}

type NamespaceQuotaGroup struct {
	BaseModel
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	ProjectQuotaGroupID uuid.UUID          `gorm:"type:uuid;not null" json:"project_quota_group_id"`
	Namespaces          []Namespace        `gorm:"foreignKey:QuotaGroupID" json:"namespaces"`
	Resources           []ResourceQuantity `gorm:"foreignKey:NamespaceQuotaGroupID" json:"resources"`
	ProjectID           uuid.UUID          `gorm:"type:uuid;not null" json:"project_id"`
	Project             Project            `gorm:"foreignKey:ProjectID" json:"project"`
}

type ResourceQuantity struct {
	BaseModel
	OrganizationQuotaGroupID *uuid.UUID       `gorm:"type:uuid;index" json:"organization_quota_group_id,omitempty"`
	ProjectQuotaGroupID      *uuid.UUID       `gorm:"type:uuid;index" json:"project_quota_group_id,omitempty"`
	NamespaceQuotaGroupID    *uuid.UUID       `gorm:"type:uuid;index" json:"namespace_quota_group_id,omitempty"`
	Quantity                 uint             `json:"quantity"`
	ResourcePropertyID       uuid.UUID        `gorm:"type:uuid;not null" json:"resource_property_id"`
	ResourceProperty         ResourceProperty `gorm:"foreignKey:ResourcePropertyID" json:"resource_property"`
}

type ResourceProperty struct {
	BaseModel
	ResourceID uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	Resource   Resource  `gorm:"foreignKey:ResourceID" json:"resource"`
	Price      float32   `gorm:"not null" json:"price"`
	Duration   float32   `json:"duration"`
}

// Refactor

// type ResourceQuantity struct {
// 	BaseModel
// 	Level          string           `gorm:"type:varchar(20);not null"` // "org", "project", "namespace"
// 	QuotaID        uuid.UUID        `gorm:"type:uuid;index;not null"`  // link to OrganizationQuota / ProjectQuota / NamespaceQuota
// 	Quantity       uint             `json:"quantity"`
// 	ResourcePropID uuid.UUID        `gorm:"type:uuid;not null" json:"resource_property_id"`
// 	ResourceProp   ResourceProperty `gorm:"foreignKey:ResourcePropID" json:"resource_property"`
// }

type OrganizationQuota struct {
	BaseModel
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	ResourcePoolID uuid.UUID          `gorm:"type:uuid;not null" json:"resource_pool_id"`
	FromOrgID      uuid.UUID          `gorm:"type:uuid;not null" json:"from_organization_id"`
	ToOrgID        uuid.UUID          `gorm:"type:uuid;not null" json:"to_organization_id"`
	FromOrg        Organization       `gorm:"foreignKey:FromOrgID" json:"from_organization"`
	ToOrg          Organization       `gorm:"foreignKey:ToOrgID" json:"to_organization"`
	Resources      []ResourceQuantity `gorm:"foreignKey:GroupID" json:"resources"`
}

type ProjectQuota struct {
	BaseModel
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	OrganizationQuotaID uuid.UUID          `gorm:"type:uuid;not null" json:"organization_quota_id"`
	ProjectID           uuid.UUID          `gorm:"type:uuid;not null" json:"project_id"`
	ResourcePoolID      uuid.UUID          `gorm:"type:uuid;not null" json:"resource_pool_id"`
	Project             Project            `gorm:"foreignKey:ProjectID" json:"project"`
	Resources           []ResourceQuantity `gorm:"foreignKey:GroupID" json:"resources"`
}

type NamespaceQuota struct {
	BaseModel
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	ProjectQuotaID uuid.UUID          `gorm:"type:uuid;not null" json:"project_quota_id"`
	NamespaceID    uuid.UUID          `gorm:"type:uuid;not null" json:"namespace_id"`
	ResourcePoolID uuid.UUID          `gorm:"type:uuid;not null" json:"resource_pool_id"`
	Namespaces     []Namespace        `gorm:"many2many:namespace_quotas;" json:"namespaces"`
	Resources      []ResourceQuantity `gorm:"foreignKey:GroupID" json:"resources"`
}

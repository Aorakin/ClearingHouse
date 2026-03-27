package models

import "github.com/google/uuid"

type Resource struct {
	BaseModel
	Name           string       `json:"name"`
	Quantity       uint         `json:"quantity"`
	ResourceTypeID uuid.UUID    `gorm:"type:uuid;not null" json:"resource_type_id"`
	ResourceType   ResourceType `gorm:"foreignKey:ResourceTypeID" json:"resource_type"`
	NodeID         uuid.UUID    `gorm:"type:uuid" json:"node_id,omitempty"`
	Node           ResourceNode `gorm:"foreignKey:NodeID" json:"-"`
	// ResourcePoolID uuid.UUID    `gorm:"type:uuid;not null" json:"resource_pool_id"`
	// ResourcePool   ResourcePool `gorm:"foreignKey:ResourcePoolID" json:"-"`
}

type ResourceType struct {
	BaseModel
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type ResourcePool struct {
	BaseModel
	Name           string         `gorm:"not null" json:"name"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index" json:"organization_id"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID" json:"-"`
	GlideletURN    string         `json:"glidelet_urn"`
	Nodes          []ResourceNode `gorm:"foreignKey:ResourcePoolID" json:"nodes"`
}

type ResourceNode struct {
	BaseModel
	Name           string       `gorm:"not null" json:"name"`
	DisplayName    string       `json:"display_name"`
	ResourcePoolID uuid.UUID    `gorm:"type:uuid;not null;index" json:"resource_pool_id"`
	ResourcePool   ResourcePool `gorm:"foreignKey:ResourcePoolID" json:"-"`
	Resources      []Resource   `gorm:"foreignKey:NodeID" json:"resources"`
}

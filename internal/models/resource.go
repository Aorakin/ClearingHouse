package models

import "github.com/google/uuid"

type Resource struct {
	BaseModel
	Name           string       `json:"name"`
	Quantity       uint         `json:"quantity"`
	ResourceTypeID uuid.UUID    `gorm:"type:uuid;not null" json:"resource_type_id"`
	ResourceType   ResourceType `gorm:"foreignKey:ResourceTypeID" json:"resource_type"`
	ResourcePoolID uuid.UUID    `gorm:"type:uuid;not null" json:"resource_pool_id"`
	ResourcePool   ResourcePool `gorm:"foreignKey:ResourcePoolID" json:"resource_pool"`
}

type ResourceType struct {
	BaseModel
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type ResourcePool struct {
	BaseModel
	Name           string       `json:"name"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Resources      []Resource   `gorm:"foreignKey:ResourcePoolID" json:"resources"`
}

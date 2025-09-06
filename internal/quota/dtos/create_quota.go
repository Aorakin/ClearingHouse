package dtos

import "github.com/google/uuid"

type CreateOrganizationQuotaRequest struct {
	Name               string                  `json:"name" binding:"required"`
	Description        string                  `json:"description"`
	FromOrganizationID uuid.UUID               `json:"from_organization_id" binding:"required,uuid"`
	ToOrganizationID   uuid.UUID               `json:"to_organization_id" binding:"required,uuid"`
	Resources          []OrganizationResources `json:"resources" binding:"required"`
}

type OrganizationResources struct {
	Quantity   uint      `json:"quantity" binding:"required"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
	Price      float32   `json:"price" binding:"required,gte=0"`
	Duration   float32   `json:"duration" binding:"gte=0"`
}

type CreateQuotaRequest struct {
	Name           string     `json:"name" binding:"required"`
	Description    string     `json:"description"`
	OrganizationID uuid.UUID  `json:"organization_id" binding:"required,uuid"`
	Resources      []Resource `json:"resources" binding:"required"`
}

type Resource struct {
	Quantity   uint      `json:"quantity" binding:"required"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
}

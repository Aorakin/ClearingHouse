package dtos

import "github.com/google/uuid"

type CreateOrganizationQuotaRequest struct {
	Name               string                     `json:"name" binding:"required"`
	Description        string                     `json:"description"`
	FromOrganizationID uuid.UUID                  `json:"from_organization_id" binding:"required,uuid"`
	ToOrganizationID   uuid.UUID                  `json:"to_organization_id" binding:"required,uuid"`
	Resources          []OrganizationQuotaRequest `json:"resources" binding:"required"`
}

type OrganizationQuotaRequest struct {
	Quantity   uint      `json:"quantity" binding:"required"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
	Price      float32   `json:"price" binding:"required"`
}

type CreateProjectQuotaRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateNamespaceQuotaRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

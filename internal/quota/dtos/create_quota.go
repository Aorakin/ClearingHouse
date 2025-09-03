package dtos

import "github.com/google/uuid"

type CreateOrganizationQuotaRequest struct {
	Name               string                     `json:"name" binding:"required"`
	Description        string                     `json:"description"`
	FromOrganizationID uuid.UUID                  `json:"from_organization_id" binding:"required,uuid"`
	ToOrganizationID   uuid.UUID                  `json:"to_organization_id" binding:"required,uuid"`
	Resources          []OrganizationQuotaRequest `json:"resources" binding:"required"`
	Creator            uuid.UUID
}

type OrganizationQuotaRequest struct {
	Quantity   uint      `json:"quantity" binding:"required"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
	Price      float32   `json:"price" binding:"required"`
}

type CreateProjectQuotaRequest struct {
	Name           string         `json:"name" binding:"required"`
	Description    string         `json:"description"`
	OrganizationID uuid.UUID      `json:"organization_id" binding:"required,uuid"`
	ResourcePools  []ResourcePool `json:"resource_pools" binding:"required"`
	Creator        uuid.UUID
}

type CreateNamespaceQuotaRequest struct {
	Name          string         `json:"name" binding:"required"`
	Description   string         `json:"description"`
	ProjectID     uuid.UUID      `json:"project_id" binding:"required,uuid"`
	ResourcePools []ResourcePool `json:"resource_pools" binding:"required"`
	Creator       uuid.UUID
}

type ResourcePool struct {
	QuotaGroupID uuid.UUID   `json:"quota_group_id" binding:"required,uuid"`
	Resources    []Resources `json:"resources" binding:"required"`
}

type Resources struct {
	Quantity   uint      `json:"quantity" binding:"required,gte=0"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
}

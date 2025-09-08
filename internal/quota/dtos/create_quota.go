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

type CreateProjectQuotaRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	ProjectID   uuid.UUID  `json:"project_id" binding:"required,uuid"`
	OrgQuotaID  uuid.UUID  `json:"org_quota_id" binding:"required,uuid"`
	OrgID       uuid.UUID  `json:"org_id" binding:"required,uuid"`
	Resources   []Resource `json:"resources" binding:"required"`
}

type CreateOwnedProjectQuotaRequest struct {
	Name           string     `json:"name" binding:"required"`
	Description    string     `json:"description"`
	ProjectID      uuid.UUID  `json:"project_id" binding:"required,uuid"`
	OrgID          uuid.UUID  `json:"org_id" binding:"required,uuid"`
	ResourcePoolID uuid.UUID  `json:"resource_pool_id" binding:"required,uuid"`
	Resources      []Resource `json:"resources" binding:"required"`
}

type CreateNamespaceQuotaRequest struct {
	Name           string     `json:"name" binding:"required"`
	Description    string     `json:"description"`
	ProjectID      uuid.UUID  `json:"project_id" binding:"required,uuid"`
	ProjectQuotaID uuid.UUID  `json:"project_quota_id" binding:"required,uuid"`
	Resources      []Resource `json:"resources" binding:"required"`
}

type Resource struct {
	Quantity   uint      `json:"quantity" binding:"required"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
}

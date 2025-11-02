package dtos

import "github.com/google/uuid"

type FindOrganizationQuotaGroupRequest struct {
	FromOrganizationID string `form:"from" binding:"required,uuid"`
	ToOrganizationID   string `form:"to" binding:"required,uuid"`
}

type NamespaceQuotaResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	ResourcePoolID   uuid.UUID `json:"resource_pool_id"`
	ResourcePoolName string    `json:"resource_pool_name"`
	OrganizationName string    `json:"organization_name"`
	NamespaceID      uuid.UUID `json:"namespace_id"`
}

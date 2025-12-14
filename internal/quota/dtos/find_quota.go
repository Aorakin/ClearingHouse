package dtos

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type FindOrganizationQuotaGroupRequest struct {
	FromOrganizationID string `form:"from" binding:"required,uuid"`
	ToOrganizationID   string `form:"to" binding:"required,uuid"`
}

type NamespaceQuotaResponse struct {
	ID               uuid.UUID                 `json:"id"`
	Name             string                    `json:"name"`
	NodeID           uuid.UUID                 `json:"node_id"`
	NodeName         string                    `json:"node_name"`
	OrganizationName string                    `json:"organization_name"`
	ProjectID        uuid.UUID                 `json:"project_id"`
	Resources        []models.ResourceQuantity `json:"resources"`
}

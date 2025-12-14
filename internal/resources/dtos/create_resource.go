package dtos

import (
	"github.com/google/uuid"
)

type UpdateResourceRequest struct {
	Quantity uint   `json:"quantity" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type CreateResourcePoolRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" binding:"required,uuid"`
	Name           string    `json:"name" binding:"required"`
}

type CreateResourceNodeRequest struct {
	ResourcePoolID uuid.UUID `json:"resource_pool_id" binding:"required,uuid"`
	NodeName       string    `json:"node_name" binding:"required"`
}

type CreateResourceTypeRequest struct {
	Unit string `json:"unit" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type CreateResourceRequest struct {
	ResourceNodeID uuid.UUID `json:"resource_node_id" binding:"required,uuid"`
	ResourceTypeID uuid.UUID `json:"resource_type_id" binding:"required,uuid"`
	Quantity       uint      `json:"quantity" binding:"required"`
	Name           string    `json:"name" binding:"required"`
}

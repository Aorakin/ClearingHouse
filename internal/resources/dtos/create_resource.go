package dtos

import (
	"github.com/google/uuid"
)

type CreateResourcePoolRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" binding:"required,uuid"`
	Name           string    `json:"name" binding:"required"`
}

type CreateResourceRequest struct {
	ResourcePoolID uuid.UUID `json:"resource_pool_id" binding:"required,uuid"`
	ResourceTypeID uuid.UUID `json:"resource_type_id" binding:"required,uuid"`
	Quantity       uint      `json:"quantity" binding:"required"`
	Name           string    `json:"name" binding:"required"`
}

type UpdateResourceRequest struct {
	Quantity uint   `json:"quantity" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type CreateResourceTypeRequest struct {
	Unit string `json:"unit" binding:"required"`
	Name string `json:"name" binding:"required"`
}

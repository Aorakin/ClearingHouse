package dtos

import "github.com/google/uuid"

type CreateNamespaceQuotaRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	NamespaceID uuid.UUID  `json:"namespace_id" binding:"required,uuid"`
	PoolID      uuid.UUID  `json:"pool_id" binding:"required,uuid"`
	Resources   []Resource `json:"resources" binding:"required"`
}

type Resource struct {
	Quantity   uint      `json:"quantity" binding:"required"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
	Price      float32   `json:"price" binding:"required,gte=0"`
	Duration   uint      `json:"duration" binding:"gte=1"`
}

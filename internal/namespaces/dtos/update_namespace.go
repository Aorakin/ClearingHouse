package dtos

import "github.com/google/uuid"

type UpdateNamespaceRequest struct {
	NamespaceID uuid.UUID `json:"namespace_id" binding:"required"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Credit      *float32  `json:"credit" binding:"omitempty,gte=0"`
}

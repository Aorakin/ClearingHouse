package dtos

import "github.com/google/uuid"

type CreatePrivateNamespaceRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required,uuid"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Credit      float32   `json:"credit" binding:"required,gte=0"`
}

package dtos

import "github.com/google/uuid"

type CreateNamespaceRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description" binding:"required"`
	Credit      float32     `json:"credit" binding:"required,gte=0"`
	ProjectID   uuid.UUID   `json:"project_id" binding:"required,uuid"`
	Members     []uuid.UUID `json:"members" binding:"dive,uuid"`
}

package dtos

import "github.com/google/uuid"

type CreateNamespaceQuotaTemplateRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	ProjectID   uuid.UUID   `json:"project_id" binding:"required,uuid"`
	QuotaIDs    []uuid.UUID `json:"quota_ids" binding:"required"`
}

package dtos

import "github.com/google/uuid"

type AssignQuotaToNamespaceRequest struct {
	Namespaces      uuid.UUIDs `json:"namespaces" binding:"required,dive,uuid"`
	ProjectID       uuid.UUID  `json:"project_id" binding:"required,uuid"`
	QuotaTemplateID uuid.UUID  `json:"quota_template_id" binding:"required"`
}

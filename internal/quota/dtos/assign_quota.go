package dtos

import "github.com/google/uuid"

type AssignQuotaToNamespaceRequest struct {
	Namespaces uuid.UUIDs `json:"namespaces" binding:"required,dive,uuid"`
	ProjectID  uuid.UUID  `json:"project_id" binding:"required,uuid"`
	QuotaID    uuid.UUID  `json:"quota_id" binding:"required"`
}

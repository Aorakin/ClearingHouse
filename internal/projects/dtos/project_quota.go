package dtos

import "github.com/google/uuid"

type ProjectQuotaResponse struct {
	ProjectID uuid.UUID `json:"project_id"`
	Quota     int       `json:"quota"`
}

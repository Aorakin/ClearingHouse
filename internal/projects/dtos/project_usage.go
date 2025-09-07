package dtos

import "github.com/google/uuid"

type ProjectUsageResponse struct {
	ProjectID uuid.UUID `json:"project_id"`
	Usage     int       `json:"usage"`
}

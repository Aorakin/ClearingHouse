package dtos

import "github.com/google/uuid"

type CreateProjectRequest struct {
	Name           string    `json:"name" binding:"required"`
	Description    string    `json:"description"`
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
}

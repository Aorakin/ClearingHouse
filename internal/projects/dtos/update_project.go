package dtos

import "github.com/google/uuid"

type UpdateProjectRequest struct {
	ProjectID   uuid.UUID `json:"project_id" binding:"required"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

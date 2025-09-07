package dtos

import "github.com/google/uuid"

type AddMembersRequest struct {
	ProjectID uuid.UUID   `json:"project_id" binding:"required,uuid"`
	Members   []uuid.UUID `json:"members" binding:"required,dive,uuid"`
}

package dtos

import "github.com/google/uuid"

type AddAdminsRequest struct {
	ProjectID uuid.UUID   `json:"project_id" binding:"required,uuid"`
	Admins    []uuid.UUID `json:"admins" binding:"required,dive,uuid"`
}

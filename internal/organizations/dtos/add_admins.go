package dtos

import "github.com/google/uuid"

type AddAdminsRequest struct {
	OrganizationID uuid.UUID   `json:"organization_id" binding:"required,uuid"`
	Admins         []uuid.UUID `json:"admins" binding:"required,dive,uuid"`
}

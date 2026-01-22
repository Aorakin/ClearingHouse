package dtos

import "github.com/google/uuid"

type RemoveMembersRequest struct {
	OrganizationID uuid.UUID   `json:"organization_id" binding:"required,uuid"`
	Members        []uuid.UUID `json:"members" binding:"required,dive,uuid"`
}

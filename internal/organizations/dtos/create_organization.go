package dtos

import "github.com/google/uuid"

type CreateOrganization struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Creator     uuid.UUID `json:"creator"`
}

type OrganizationURI struct {
	OrgID     string `uri:"id" binding:"required,uuid"`
	RequestID uuid.UUID
}

type UpdateOrganization struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

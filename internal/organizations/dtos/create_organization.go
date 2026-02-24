package dtos

import "github.com/google/uuid"

type CreateOrganization struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
}

type OrganizationURI struct {
	OrgID     string `uri:"org-id" binding:"required,uuid"`
	RequestID uuid.UUID
}

type UpdateOrganization struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Domain      string `json:"domain"`
}

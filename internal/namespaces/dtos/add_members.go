package dtos

import (
	"github.com/google/uuid"
)

type AddMembersRequest struct {
	NamespaceID uuid.UUID   `json:"namespace_id" binding:"required,uuid"`
	Members     []uuid.UUID `json:"members" binding:"required,dive,uuid"`
	Creator     uuid.UUID
}

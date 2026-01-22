package dtos

import (
	"github.com/google/uuid"
)

type RemoveMembersRequest struct {
	NamespaceID uuid.UUID   `json:"namespace_id" binding:"required,uuid"`
	Members     []uuid.UUID `json:"members" binding:"required,dive,uuid"`
}

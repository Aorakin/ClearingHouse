package dtos

import "github.com/google/uuid"

type CreateTicketRequest struct {
	Name        string      `json:"name" binding:"required"`
	NamespaceID uuid.UUID   `json:"namespace_id" binding:"required,uuid"`
	QuotaID     uuid.UUID   `json:"quota_id" binding:"required,uuid"`
	Resources   []Resources `json:"resources" binding:"required"`
	Duration    float32     `json:"duration" binding:"required,gte=1"`
}

type Resources struct {
	Quantity   uint      `json:"quantity" binding:"required,gte=0"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
}

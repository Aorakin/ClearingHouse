package dtos

import "github.com/google/uuid"

type CreateTicketRequest struct {
	Name        string      `json:"name" binding:"required"`
	NamespaceID uuid.UUID   `json:"namespace_id" binding:"required,uuid"`
	Creator     uuid.UUID   `json:"creator"`
	Resources   []Resources `json:"resources" binding:"required"`
	Duration    uint        `json:"duration" binding:"required,gte=1"`
}

type Resources struct {
	Quantity   uint      `json:"quantity" binding:"required,gte=0"`
	ResourceID uuid.UUID `json:"resource_id" binding:"required,uuid"`
}

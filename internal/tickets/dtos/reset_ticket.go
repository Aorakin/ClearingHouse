package dtos

import "github.com/google/uuid"

type ResetTicketsRequest struct {
	TicketIDs []uuid.UUID `json:"ticket_ids" binding:"required,dive,uuid"`
}

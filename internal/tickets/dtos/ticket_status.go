package dtos

import (
	"time"

	"github.com/google/uuid"
)

type StartTicketsRequest struct {
	Tickets []StartTicketRequest `json:"tickets" binding:"required,dive"`
}

type StartTicketRequest struct {
	TicketID  uuid.UUID `json:"ticket_id" binding:"required,uuid"`
	StartTime time.Time `json:"start_time" binding:"required"`
}

type StopTicketsRequest struct {
	Tickets []StopTicketRequest `json:"tickets" binding:"required,dive"`
}

type StopTicketRequest struct {
	TicketID uuid.UUID `json:"ticket_id" binding:"required,uuid"`
}

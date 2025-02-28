package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/dto"
	"github.com/google/uuid"
)

type TicketsUsecase interface {
	Create(namespaceID uuid.UUID, ticketRequest *dto.CreateTicketRequest) (*models.Ticket, error)
}

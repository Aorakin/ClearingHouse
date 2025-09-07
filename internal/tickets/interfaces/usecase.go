package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/google/uuid"
)

type TicketUsecase interface {
	CreateTicket(request *dtos.CreateTicketRequest, userID uuid.UUID) (*models.Ticket, error)
	GetNamespaceTickets(namespaceID uuid.UUID, userID uuid.UUID) ([]models.Ticket, error)
}

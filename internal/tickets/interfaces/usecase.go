package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/google/uuid"
)

type TicketUsecase interface {
	CreateTicket(request *dtos.CreateTicketRequest, userID uuid.UUID) (*dtos.GliderTicketResponse, error)
	GetNamespaceTickets(namespaceID uuid.UUID, userID uuid.UUID) ([]models.Ticket, error)
	GetUserTickets(userID uuid.UUID) ([]dtos.TicketResponse, error)
	StartTicket(request *dtos.StartTicketsRequest) ([]models.Ticket, error)
	StopTicket(request *dtos.StopTicketsRequest) ([]models.Ticket, error)
	GetTicket(ticketID uuid.UUID, userID uuid.UUID) (*models.Ticket, error)
	CancelTicket(ticketID uuid.UUID, userID uuid.UUID) error
}

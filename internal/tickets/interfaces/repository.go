package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type TicketRepository interface {
	GetTicketByID(ticketID uuid.UUID) (*models.Ticket, error)
	CreateTicket(ticket *models.Ticket) error
	CreateTicketResource(resource *models.TicketResource) error
	GetNamespaceUsage(namespaceID, quotaID, resourceID uuid.UUID) (uint, error)

	GetResourceUsage(namespaceID, quotaID, resourceID uuid.UUID) (uint, error)

	GetNamespaceTickets(namespaceID uuid.UUID) ([]models.Ticket, error)
}

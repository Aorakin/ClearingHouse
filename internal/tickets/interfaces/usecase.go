package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/dtos"
)

type TicketUsecase interface {
	CreateTicket(request *dtos.CreateTicketRequest) (*models.Ticket, error)
}

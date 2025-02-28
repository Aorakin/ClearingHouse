package usecase

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/dto"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	"github.com/google/uuid"
)

type TicketsUsecase struct {
	TicketsRepository interfaces.TicketsRepository
}

func NewTicketsUsecase(TicketsRepository interfaces.TicketsRepository) interfaces.TicketsUsecase {
	return &TicketsUsecase{
		TicketsRepository: TicketsRepository,
	}
}

func (u *TicketsUsecase) Create(namespaceID uuid.UUID, ticketRequest *dto.CreateTicketRequest) (*models.Ticket, error) {
	return nil, nil
}

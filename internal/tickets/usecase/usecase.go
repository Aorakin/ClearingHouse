package usecase

import "github.com/ClearingHouse/internal/tickets/interfaces"

type TicketsUsecase struct {
	TicketsRepository interfaces.TicketsRepository
}

func NewTicketsUsecase(TicketsRepository interfaces.TicketsRepository) interfaces.TicketsUsecase {
	return &TicketsUsecase{
		TicketsRepository: TicketsRepository,
	}
}

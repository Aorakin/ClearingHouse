package usecase

import "github.com/ClearingHouse/internal/quotas/interfaces"

type QuotasUsecase struct {
	quotasRepository interfaces.QuotasRepository
}

func NewQuotasUsecase(quotasRepository interfaces.QuotasRepository) interfaces.QuotasUsecase {
	return &QuotasUsecase{quotasRepository: quotasRepository}
}


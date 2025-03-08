package usecase

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quotas/dto"
	"github.com/ClearingHouse/internal/quotas/interfaces"
	"github.com/google/uuid"
)

type QuotasUsecase struct {
	quotasRepository interfaces.QuotasRepository
}

func NewQuotasUsecase(quotasRepository interfaces.QuotasRepository) interfaces.QuotasUsecase {
	return &QuotasUsecase{quotasRepository: quotasRepository}
}

func (u *QuotasUsecase) GetAll(userID uuid.UUID, projectID uuid.UUID) ([]models.Quota, error) {
	var quotas []models.Quota
	return quotas, nil
}

func (u *QuotasUsecase) Create(userID uuid.UUID, request *dto.CreateQuotaRequest) (*models.Quota, error) {
	var quota models.Quota
	return &quota, nil
}

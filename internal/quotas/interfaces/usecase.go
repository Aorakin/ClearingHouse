package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quotas/dto"
	"github.com/google/uuid"
)

type QuotasUsecase interface {
	GetAll(userID uuid.UUID, projectID uuid.UUID) ([]models.Quota, error)
	Create(userID uuid.UUID, request *dto.CreateQuotaRequest) (*models.Quota, error)
	// Delete(userID uuid.UUID, quotaID uuid.UUID) error
}

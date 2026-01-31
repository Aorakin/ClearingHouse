package usecase

import (
	"errors"
	"fmt"

	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) DeleteProjectQuota(quotaID uuid.UUID, userID uuid.UUID) error {
	quota, err := u.quotaRepo.GetProjectQuotaByID(quotaID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("project quota not found: %w", err))
	}

	if err := u.isOrgAdmin(quota.OrganizationID, userID); err != nil {
		return err
	}

	hasNamespaceQuotas, err := u.quotaRepo.HasNamespaceQuotasByProjectQuotaID(quotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to check namespace quotas: %w", err))
	}
	if hasNamespaceQuotas {
		return apiError.NewBadRequestError(errors.New("cannot delete project quota: namespace quotas are still using this project quota"))
	}

	if err := u.quotaRepo.DeleteProjectQuota(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota: %w", err))
	}

	return nil
}

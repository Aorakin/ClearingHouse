package usecase

import (
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

	return u.deleteProjectQuotaCascade(quotaID)
}

// deleteProjectQuotaCascade cascades deletion down through every NamespaceQuota
// that belongs to this ProjectQuota, then deletes the project quota itself.
func (u *QuotaUsecase) deleteProjectQuotaCascade(quotaID uuid.UUID) error {
	nsQuotas, err := u.quotaRepo.GetNamespaceQuotasByProjectQuotaID(quotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to list namespace quotas: %w", err))
	}

	for _, nsq := range nsQuotas {
		if err := u.deleteNamespaceQuotaCascade(nsq.ID); err != nil {
			return err
		}
	}

	if err := u.quotaRepo.DeleteResourceQuantitiesByProjectQuotaID(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete resource quantities for project quota: %w", err))
	}

	if err := u.quotaRepo.DeleteProjectQuota(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota: %w", err))
	}

	return nil
}

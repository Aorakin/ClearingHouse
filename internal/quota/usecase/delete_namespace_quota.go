package usecase

import (
	"fmt"

	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) DeleteNamespaceQuota(quotaID uuid.UUID, userID uuid.UUID) error {
	quota, err := u.quotaRepo.GetNamespaceQuotaByID(quotaID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("namespace quota not found: %w", err))
	}

	if quota.ProjectID == nil {
		return apiError.NewInternalServerError(fmt.Errorf("namespace quota has no associated project"))
	}

	if err := u.isProjAdmin(*quota.ProjectID, userID); err != nil {
		return err
	}

	return u.deleteNamespaceQuotaCascade(quotaID)
}

// deleteNamespaceQuotaCascade removes all junction-table references from templates,
// deletes resource quantities, then soft-deletes the namespace quota itself.
// It does NOT delete templates — the many2many relationship means templates may
// reference other quotas that should be preserved.
func (u *QuotaUsecase) deleteNamespaceQuotaCascade(quotaID uuid.UUID) error {
	if err := u.quotaRepo.RemoveNamespaceQuotaFromAllTemplates(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to remove namespace quota from templates: %w", err))
	}

	if err := u.quotaRepo.DeleteResourceQuantitiesByNamespaceQuotaID(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete resource quantities: %w", err))
	}

	if err := u.quotaRepo.DeleteNamespaceQuota(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete namespace quota: %w", err))
	}

	return nil
}

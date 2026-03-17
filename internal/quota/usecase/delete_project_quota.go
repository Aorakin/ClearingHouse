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

	// Cascading soft-delete: project quota -> namespace quotas -> resource quantities
	nsQuotaIDs, err := u.quotaRepo.SoftDeleteNamespaceQuotasByProjectQuotaID(quotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to cascade delete namespace quotas: %w", err))
	}

	// Clean up template associations for deleted namespace quotas
	if err := u.quotaRepo.UnassignQuotaTemplatesByNamespaceQuotaIDs(nsQuotaIDs); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to clean up quota template associations: %w", err))
	}

	// Soft-delete resource quantities for each namespace quota
	for _, nsID := range nsQuotaIDs {
		if err := u.quotaRepo.DeleteResourceQuantitiesByNamespaceQuotaID(nsID); err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to delete namespace quota resource quantities: %w", err))
		}
	}

	// Soft-delete resource quantities for the project quota
	if err := u.quotaRepo.SoftDeleteResourceQuantitiesByProjectQuotaID(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota resource quantities: %w", err))
	}

	// Finally, soft-delete the project quota itself
	if err := u.quotaRepo.DeleteProjectQuota(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota: %w", err))
	}

	return nil
}

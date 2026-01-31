package usecase

import (
	"errors"
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
		return apiError.NewInternalServerError(errors.New("namespace quota has no associated project"))
	}

	if err := u.isProjAdmin(*quota.ProjectID, userID); err != nil {
		return err
	}

	hasQuotaTemplates, err := u.quotaRepo.HasQuotaTemplatesByNamespaceQuotaID(quotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to check quota templates: %w", err))
	}
	if hasQuotaTemplates {
		return apiError.NewBadRequestError(errors.New("cannot delete namespace quota: this quota is being used in quota templates"))
	}

	if err := u.quotaRepo.DeleteNamespaceQuota(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete namespace quota: %w", err))
	}

	return nil
}

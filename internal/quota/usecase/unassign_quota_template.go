package usecase

import (
	"fmt"

	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) UnassignQuotaTemplateFromNamespace(namespaceID uuid.UUID, userID uuid.UUID) error {
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("namespace not found: %w", err))
	}

	if namespace.ProjectID == nil {
		return apiError.NewBadRequestError("cannot unassign quota template from private namespace")
	}

	if err := u.isProjAdmin(*namespace.ProjectID, userID); err != nil {
		return err
	}

	if namespace.QuotaTemplateID == nil {
		return apiError.NewBadRequestError("namespace does not have a quota template assigned")
	}

	if err := u.quotaRepo.UnassignQuotaTemplateFromNamespace(namespaceID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to unassign quota template: %w", err))
	}

	return nil
}

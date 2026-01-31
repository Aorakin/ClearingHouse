package usecase

import (
	"errors"
	"fmt"

	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) DeleteNamespaceQuotaTemplate(templateID uuid.UUID, userID uuid.UUID) error {
	template, err := u.quotaRepo.GetNamespaceQuotaTemplateByID(templateID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("namespace quota template not found: %w", err))
	}

	if err := u.isProjAdmin(template.ProjectID, userID); err != nil {
		return err
	}

	hasNamespaces, err := u.quotaRepo.HasNamespacesUsingTemplate(templateID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to check namespaces using template: %w", err))
	}
	if hasNamespaces {
		return apiError.NewBadRequestError(errors.New("cannot delete quota template: namespaces are still using this template"))
	}

	if err := u.quotaRepo.DeleteNamespaceQuotaTemplate(templateID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete quota template: %w", err))
	}

	return nil
}

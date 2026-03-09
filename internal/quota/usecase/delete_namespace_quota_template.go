package usecase

import (
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

	// Unassign every namespace that references this template.
	// Namespaces are not deleted — they simply lose their quota assignment.
	if err := u.quotaRepo.UnassignAllNamespacesFromTemplate(templateID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to unassign namespaces from template: %w", err))
	}

	// Remove all quota-template relations from the junction table.
	// The NamespaceQuotas themselves are preserved; they may belong to other templates.
	if err := u.quotaRepo.DeleteTemplateQuotaRelations(templateID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete template quota relations: %w", err))
	}

	if err := u.quotaRepo.DeleteNamespaceQuotaTemplate(templateID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete quota template: %w", err))
	}

	return nil
}

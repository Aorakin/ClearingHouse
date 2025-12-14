package usecase

import (
	"errors"
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) CreateProjectQuota(request *dtos.CreateProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error) {
	if err := u.isOrgAdmin(request.OrgID, userID); err != nil {
		return nil, err
	}

	quota, err := u.quotaRepo.GetOrgQuotaByID(request.OrgQuotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find organization quota: %w", err))
	}

	if quota.ToOrgID != request.OrgID {
		return nil, apiError.NewForbiddenError(errors.New("organization quota does not belong to the organization"))
	}

	if quota.NodeID != request.NodeID {
		return nil, apiError.NewForbiddenError(errors.New("resource node does not belong to the organization quota"))
	}

	quotaResourcesMap := make(map[uuid.UUID]models.ResourceQuantity)
	for _, resource := range quota.Resources {
		quotaResourcesMap[resource.ResourceProp.ResourceID] = resource
	}

	if err := u.validateProjectQuotaRequest(request, quotaResourcesMap); err != nil {
		return nil, err
	}

	return u.createProjectQuota(request, quotaResourcesMap)
}

func (u *QuotaUsecase) GetProjectQuotas(projectID uuid.UUID) ([]models.ProjectQuota, error) {
	quotas, err := u.quotaRepo.GetProjectQuotaByProjectID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get project quotas: %w", err))
	}
	return quotas, nil
}

func (u *QuotaUsecase) validateProjectQuotaRequest(request *dtos.CreateProjectQuotaRequest, quotaResourcesMap map[uuid.UUID]models.ResourceQuantity) error {
	quotaExists, err := u.quotaRepo.IsProjectQuotaExist(request.ProjectID, request.NodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to check existing project quota: %w", err))
	}
	if quotaExists {
		return apiError.NewConflictError(errors.New("project quota already exists for this project and node"))
	}

	quota, err := u.quotaRepo.GetOrgQuotaByID(request.OrgQuotaID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to find organization quota: %w", err))
	}

	seenResources := make(map[uuid.UUID]struct{})
	for _, r := range request.Resources {
		if r.Quantity <= 0 {
			return apiError.NewBadRequestError(fmt.Errorf("quantity for resource %s must be greater than zero", r.ResourceID))
		}

		if _, exists := quotaResourcesMap[r.ResourceID]; !exists {
			return apiError.NewBadRequestError(fmt.Errorf("resource %s not found in organization quota", r.ResourceID))
		}

		if _, duplicate := seenResources[r.ResourceID]; duplicate {
			return apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
		}
		seenResources[r.ResourceID] = struct{}{}

		currentUsage, err := u.quotaRepo.GetOrgUsage(quota.ID, r.ResourceID)
		if err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to get current usage for resource %s: %w", r.ResourceID, err))
		}

		maxQuota, err := u.quotaRepo.GetOrgQuotaQuantity(quota.ID, r.ResourceID)
		if err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to get maximum quota for resource %s: %w", r.ResourceID, err))
		}

		if r.Quantity+currentUsage > maxQuota {
			return apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quota %d for resource %s", r.Quantity, maxQuota-currentUsage, r.ResourceID))
		}
	}

	return nil
}

func (u *QuotaUsecase) createProjectQuota(request *dtos.CreateProjectQuotaRequest, quotaResourcesMap map[uuid.UUID]models.ResourceQuantity) (*models.ProjectQuota, error) {
	projectQuota := &models.ProjectQuota{
		Name:                request.Name,
		Description:         request.Description,
		OrganizationID:      request.OrgID,
		ProjectID:           request.ProjectID,
		OrganizationQuotaID: &request.OrgQuotaID,
		NodeID:              request.NodeID,
	}

	if err := u.quotaRepo.CreateProjectQuota(projectQuota); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create project quota: %w", err))
	}

	for _, r := range request.Resources {
		resourceQuantity := models.ResourceQuantity{
			ProjectQuotaID: &projectQuota.ID,
			ResourcePropID: quotaResourcesMap[r.ResourceID].ResourceProp.ID,
			Quantity:       r.Quantity,
		}

		if err := u.quotaRepo.CreateResourceQuantity(&resourceQuantity); err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
		}
	}

	return projectQuota, nil
}

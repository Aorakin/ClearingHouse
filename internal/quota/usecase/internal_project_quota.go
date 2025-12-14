package usecase

import (
	"errors"
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) CreateInternalProjectQuota(request *dtos.CreateInternalProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error) {
	if err := u.isOrgAdmin(request.OrgID, userID); err != nil {
		return nil, err
	}

	node, err := u.resourceRepo.GetResourceNodeByID(request.NodeID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find resource node: %w", err))
	}
	if node.ResourcePool.OrganizationID != request.OrgID {
		return nil, apiError.NewForbiddenError(errors.New("resource node does not belong to the organization"))
	}

	_, err = u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find project: %w", err))
	}

	resourceProperties, err := u.validateInternalProjectQuotaRequest(request, node)
	if err != nil {
		return nil, err
	}

	return u.createInternalProjectQuota(request, resourceProperties)
}

func (u *QuotaUsecase) validateInternalProjectQuotaRequest(request *dtos.CreateInternalProjectQuotaRequest, node *models.ResourceNode) ([]dtos.ResourceWithProperty, error) {
	resourcesMap := make(map[uuid.UUID]models.Resource)
	for _, resource := range node.Resources {
		resourcesMap[resource.ID] = resource
	}

	seenResources := make(map[uuid.UUID]struct{})
	var resourceProperties []dtos.ResourceWithProperty

	for _, r := range request.Resources {
		if r.Quantity <= 0 {
			return nil, apiError.NewBadRequestError(fmt.Errorf("resource quantity must be greater than zero for resource %s", r.ResourceID))
		}

		if _, duplicate := seenResources[r.ResourceID]; duplicate {
			return nil, apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
		}

		seenResources[r.ResourceID] = struct{}{}
		if _, exists := resourcesMap[r.ResourceID]; !exists {
			return nil, apiError.NewBadRequestError(fmt.Errorf("resource %s is not in the resource pool", r.ResourceID))
		}

		if r.Quantity > resourcesMap[r.ResourceID].Quantity {
			return nil, apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, resourcesMap[r.ResourceID].Quantity, r.ResourceID))
		}

		resourceProperties = append(resourceProperties, dtos.ResourceWithProperty{
			Quantity:   r.Quantity,
			ResourceID: r.ResourceID,
			Price:      r.Price,
			Duration:   r.Duration,
		})
	}

	return resourceProperties, nil
}

func (u *QuotaUsecase) createInternalProjectQuota(request *dtos.CreateInternalProjectQuotaRequest, resourceProperties []dtos.ResourceWithProperty) (*models.ProjectQuota, error) {
	resourceQuantities, err := u.createInternalProjectResource(resourceProperties)
	if err != nil {
		return nil, err
	}

	projectQuota := &models.ProjectQuota{
		Name:           request.Name,
		Description:    request.Description,
		OrganizationID: request.OrgID,
		ProjectID:      request.ProjectID,
		NodeID:         request.NodeID,
		Resources:      resourceQuantities,
	}

	if err = u.quotaRepo.CreateProjectQuota(projectQuota); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create project quota: %w", err))
	}

	return u.quotaRepo.GetProjectQuotaByID(projectQuota.ID)
}

func (u *QuotaUsecase) createInternalProjectResource(resourceProperties []dtos.ResourceWithProperty) ([]models.ResourceQuantity, error) {
	var resourceQuantities []models.ResourceQuantity

	for _, r := range resourceProperties {
		resourceProperty := models.ResourceProperty{
			ResourceID:  r.ResourceID,
			Price:       r.Price,
			MaxDuration: r.Duration,
		}

		if err := u.quotaRepo.CreateResourceProperty(&resourceProperty); err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource property: %w", err))
		}

		resourceQuantities = append(resourceQuantities, models.ResourceQuantity{
			Quantity:       r.Quantity,
			ResourcePropID: resourceProperty.ID,
		})
	}

	return resourceQuantities, nil
}

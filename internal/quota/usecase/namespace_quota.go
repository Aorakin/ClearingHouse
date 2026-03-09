package usecase

import (
	"errors"
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error) {
	if err := u.isProjAdmin(request.ProjectID, userID); err != nil {
		return nil, err
	}

	quota, err := u.quotaRepo.GetProjectQuotaByID(request.ProjectQuotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find project quota: %w", err))
	}

	if quota.ProjectID != request.ProjectID {
		return nil, apiError.NewForbiddenError(errors.New("quota does not belong to the project"))
	}

	if quota.NodeID != request.NodeID {
		return nil, apiError.NewForbiddenError(errors.New("resource node does not belong to the project quota"))
	}

	quotaResourcesMap := make(map[uuid.UUID]models.ResourceQuantity)
	for _, resource := range quota.Resources {
		quotaResourcesMap[resource.ResourceProp.ResourceID] = resource
	}

	if err := u.validateNamespaceQuotaRequest(request, quotaResourcesMap); err != nil {
		return nil, err
	}

	return u.createNamespaceQuota(request, quotaResourcesMap)
}

func (u *QuotaUsecase) GetNamespaceQuota(namespaceID uuid.UUID) ([]dtos.NamespaceQuotaResponse, error) {
	quotas, err := u.quotaRepo.GetNamespaceQuotaByNamespaceID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quota: %w", err))
	}

	var response []dtos.NamespaceQuotaResponse
	for _, quota := range quotas {
		org, err := u.resourceRepo.GetResourceNodeOrganization(quota.NodeID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get organization for resource node %s: %w", quota.NodeID, err))
		}
		response = append(response, dtos.NamespaceQuotaResponse{
			ID:               quota.ID,
			Name:             quota.Name,
			NodeID:           quota.NodeID,
			NodeName:         quota.Node.Name,
			NodeDisplayName:  quota.Node.DisplayName,
			OrganizationName: org.Name,
			ProjectID:        *quota.ProjectID,
			Resources:        quota.Resources,
		})
	}
	return response, nil
}

func (u *QuotaUsecase) UpdateNamespaceQuota(quotaID uuid.UUID, request *dtos.UpdateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error) {
	// Get existing namespace quota
	existingQuota, err := u.quotaRepo.GetNamespaceQuotaByID(quotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find namespace quota: %w", err))
	}

	// Verify user is project admin
	if existingQuota.ProjectID == nil {
		return nil, apiError.NewBadRequestError(errors.New("namespace quota has no associated project"))
	}
	if err := u.isProjAdmin(*existingQuota.ProjectID, userID); err != nil {
		return nil, err
	}

	// Get project quota to validate resources
	if existingQuota.ProjectQuotaID == nil {
		return nil, apiError.NewBadRequestError(errors.New("namespace quota has no associated project quota"))
	}
	projectQuota, err := u.quotaRepo.GetProjectQuotaByID(*existingQuota.ProjectQuotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find project quota: %w", err))
	}

	// Build map of available project quota resources
	quotaResourcesMap := make(map[uuid.UUID]models.ResourceQuantity)
	for _, resource := range projectQuota.Resources {
		quotaResourcesMap[resource.ResourceProp.ResourceID] = resource
	}

	// Update basic fields (name and description)
	if err := u.quotaRepo.UpdateNamespaceQuota(quotaID, request.Name, request.Description); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update namespace quota: %w", err))
	}

	// Update resources if provided
	if request.Resources != nil && len(request.Resources) > 0 {
		// Validate the requested resources
		if err := u.validateUpdateNamespaceQuotaResources(request.Resources, quotaResourcesMap); err != nil {
			return nil, err
		}

		// Get existing resource quantities
		existingQuantities, err := u.quotaRepo.GetResourceQuantitiesByNamespaceQuotaID(quotaID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get existing resource quantities: %w", err))
		}

		// Build map of existing quantities by resource ID
		existingQuantitiesMap := make(map[uuid.UUID]models.ResourceQuantity)
		for _, qty := range existingQuantities {
			existingQuantitiesMap[qty.ResourceProp.ResourceID] = qty
		}

		// Update existing resources or create new ones
		for _, resource := range request.Resources {
			if existingQty, exists := existingQuantitiesMap[resource.ResourceID]; exists {
				// Update existing resource quantity
				if err := u.quotaRepo.UpdateResourceQuantity(existingQty.ID, resource.Quantity); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update resource quantity: %w", err))
				}
			} else {
				// Create new resource quantity
				resourceQuantity := models.ResourceQuantity{
					NamespaceQuotaID: &quotaID,
					ResourcePropID:   quotaResourcesMap[resource.ResourceID].ResourceProp.ID,
					Quantity:         resource.Quantity,
				}
				if err := u.quotaRepo.CreateResourceQuantity(&resourceQuantity); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
				}
			}
		}

		// Delete resources that are no longer in the request
		requestedResourceIDs := make(map[uuid.UUID]bool)
		for _, resource := range request.Resources {
			requestedResourceIDs[resource.ResourceID] = true
		}
		for resourceID, existingQty := range existingQuantitiesMap {
			if !requestedResourceIDs[resourceID] {
				if err := u.quotaRepo.UpdateResourceQuantity(existingQty.ID, 0); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to remove resource quantity: %w", err))
				}
			}
		}
	}

	// Fetch and return the updated quota
	updatedQuota, err := u.quotaRepo.GetNamespaceQuotaByID(quotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get updated namespace quota: %w", err))
	}

	return updatedQuota, nil
}

func (u *QuotaUsecase) CreateNamespaceQuotaTemplate(request *dtos.CreateNamespaceQuotaTemplateRequest, userID uuid.UUID) (*models.NamespaceQuotaTemplate, error) {
	if err := u.isProjAdmin(request.ProjectID, userID); err != nil {
		return nil, err
	}

	namespaceQuotas, err := u.quotaRepo.GetNamespaceQuotasByIDs(request.QuotaIDs, request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quotas: %w", err))
	}

	template := &models.NamespaceQuotaTemplate{
		Name:        request.Name,
		Description: request.Description,
		ProjectID:   request.ProjectID,
		Quotas:      namespaceQuotas,
	}

	if err := u.quotaRepo.CreateNamespaceQuotaTemplate(template); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create namespace quota template: %w", err))
	}

	return template, nil
}

func (u *QuotaUsecase) GetNamespaceQuotaTemplate(quotaTemplateID uuid.UUID) (*models.NamespaceQuotaTemplate, error) {
	quotaTemplate, err := u.quotaRepo.GetNamespaceQuotaTemplateByID(quotaTemplateID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find quota template: %w", err))
	}

	return quotaTemplate, nil
}

func (u *QuotaUsecase) UpdateNamespaceQuotaTemplate(quotaTemplateID uuid.UUID, request *dtos.UpdateNamespaceQuotaTemplateRequest, userID uuid.UUID) (*models.NamespaceQuotaTemplate, error) {
	// Get the existing template to check project ownership
	existingTemplate, err := u.quotaRepo.GetNamespaceQuotaTemplateByID(quotaTemplateID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find quota template: %w", err))
	}

	// Verify user is project admin
	if err := u.isProjAdmin(existingTemplate.ProjectID, userID); err != nil {
		return nil, err
	}

	// Update basic fields (name and description)
	if err := u.quotaRepo.UpdateNamespaceQuotaTemplate(quotaTemplateID, request.Name, request.Description); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update namespace quota template: %w", err))
	}

	// Update quota associations if provided
	if request.QuotaIDs != nil && len(request.QuotaIDs) > 0 {
		// Validate that all quotas belong to the same project
		namespaceQuotas, err := u.quotaRepo.GetNamespaceQuotasByIDs(request.QuotaIDs, existingTemplate.ProjectID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quotas: %w", err))
		}

		// Remove old quota associations
		if len(existingTemplate.Quotas) > 0 {
			oldQuotaIDs := make([]uuid.UUID, len(existingTemplate.Quotas))
			for i, quota := range existingTemplate.Quotas {
				oldQuotaIDs[i] = quota.ID
			}
			if err := u.quotaRepo.RemoveQuotasFromTemplate(quotaTemplateID, oldQuotaIDs); err != nil {
				return nil, apiError.NewInternalServerError(fmt.Errorf("failed to remove old quotas from template: %w", err))
			}
		}

		// Add new quota associations
		newQuotaIDs := make([]uuid.UUID, len(namespaceQuotas))
		for i, quota := range namespaceQuotas {
			newQuotaIDs[i] = quota.ID
		}
		if err := u.quotaRepo.AddQuotasToTemplate(quotaTemplateID, newQuotaIDs); err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to add quotas to template: %w", err))
		}
	}

	// Fetch and return the updated template
	updatedTemplate, err := u.quotaRepo.GetNamespaceQuotaTemplateByID(quotaTemplateID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to fetch updated template: %w", err))
	}

	return updatedTemplate, nil
}

func (u *QuotaUsecase) GetNamespaceQuotaTemplatesByProjectID(projectID uuid.UUID, userID uuid.UUID) ([]models.NamespaceQuotaTemplate, error) {
	if err := u.isProjAdmin(projectID, userID); err != nil {
		return nil, err
	}

	templates, err := u.quotaRepo.GetNamespaceQuotaTemplatesByProjectID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quota templates: %w", err))
	}

	return templates, nil
}

func (u *QuotaUsecase) AssignQuotaTemplateToNamespace(request *dtos.AssignQuotaToNamespaceRequest, userID uuid.UUID) error {
	if err := u.isProjAdmin(request.ProjectID, userID); err != nil {
		return err
	}

	quotaTemplate, err := u.quotaRepo.GetNamespaceQuotaTemplateByID(request.QuotaTemplateID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to find quota template: %w", err))
	}

	if quotaTemplate.ProjectID != request.ProjectID {
		return apiError.NewForbiddenError(errors.New("quota template does not belong to the project"))
	}

	namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(request.ProjectID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to find namespaces by project ID: %w", err))
	}
	namespaceMap := make(map[uuid.UUID]struct{})
	for _, ns := range namespaces {
		namespaceMap[ns.ID] = struct{}{}
	}

	for _, namespaceID := range request.Namespaces {
		if _, ok := namespaceMap[namespaceID]; !ok {
			return apiError.NewBadRequestError(fmt.Errorf("namespace %s does not belong to the project %s", namespaceID, request.ProjectID))
		}
	}

	for _, namespaceID := range request.Namespaces {
		if err := u.quotaRepo.AssignQuotaToNamespace(namespaceID, request.QuotaTemplateID); err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to assign quota to namespace %s: %w", namespaceID, err))
		}
	}

	return nil
}

func (u *QuotaUsecase) GetNamespaceQuotaInProject(userID uuid.UUID, projectID uuid.UUID) ([]dtos.NamespaceQuotaResponse, error) {
	if err := u.isProjAdmin(projectID, userID); err != nil {
		return nil, err
	}

	quotas, err := u.quotaRepo.GetNamespaceQuotasByProjectID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quotas in project: %w", err))
	}

	var response []dtos.NamespaceQuotaResponse
	for _, quota := range quotas {
		org, err := u.resourceRepo.GetResourceNodeOrganization(quota.NodeID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get organization for resource node %s: %w", quota.NodeID, err))
		}
		response = append(response, dtos.NamespaceQuotaResponse{
			ID:               quota.ID,
			Name:             quota.Name,
			NodeID:           quota.NodeID,
			NodeName:         quota.Node.Name,
			NodeDisplayName:  quota.Node.DisplayName,
			OrganizationName: org.Name,
			ProjectID:        *quota.ProjectID,
			Resources:        quota.Resources,
		})
	}
	return response, nil
}

func (u *QuotaUsecase) validateUpdateNamespaceQuotaResources(resources []dtos.Resource, quotaResourcesMap map[uuid.UUID]models.ResourceQuantity) error {
	if len(resources) == 0 {
		return apiError.NewBadRequestError(errors.New("at least one resource quota is required"))
	}

	seenResources := make(map[uuid.UUID]struct{})
	for _, r := range resources {
		if _, exists := quotaResourcesMap[r.ResourceID]; !exists {
			// Build list of available resource IDs for better error message
			availableResources := make([]uuid.UUID, 0, len(quotaResourcesMap))
			for resourceID := range quotaResourcesMap {
				availableResources = append(availableResources, resourceID)
			}
			return apiError.NewBadRequestError(fmt.Errorf("resource %s not found in project quota. Available resources: %v", r.ResourceID, availableResources))
		}

		if _, duplicate := seenResources[r.ResourceID]; duplicate {
			return apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
		}
		seenResources[r.ResourceID] = struct{}{}

		if r.Quantity > quotaResourcesMap[r.ResourceID].Quantity {
			return apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, quotaResourcesMap[r.ResourceID].Quantity, r.ResourceID))
		}
	}

	return nil
}

func (u *QuotaUsecase) validateNamespaceQuotaRequest(request *dtos.CreateNamespaceQuotaRequest, quotaResourcesMap map[uuid.UUID]models.ResourceQuantity) error {
	if len(request.Resources) == 0 {
		return apiError.NewBadRequestError(errors.New("at least one resource quota is required"))
	}

	seenResources := make(map[uuid.UUID]struct{})
	for _, r := range request.Resources {
		if _, exists := quotaResourcesMap[r.ResourceID]; !exists {
			// Build list of available resource IDs for better error message
			availableResources := make([]uuid.UUID, 0, len(quotaResourcesMap))
			for resourceID := range quotaResourcesMap {
				availableResources = append(availableResources, resourceID)
			}
			return apiError.NewBadRequestError(fmt.Errorf("resource %s not found in project quota. Available resources: %v", r.ResourceID, availableResources))
		}

		if _, duplicate := seenResources[r.ResourceID]; duplicate {
			return apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
		}
		seenResources[r.ResourceID] = struct{}{}

		if r.Quantity > quotaResourcesMap[r.ResourceID].Quantity {
			return apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, quotaResourcesMap[r.ResourceID].Quantity, r.ResourceID))
		}
	}

	return nil
}

func (u *QuotaUsecase) createNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, quotaResourcesMap map[uuid.UUID]models.ResourceQuantity) (*models.NamespaceQuota, error) {
	namespaceQuota := &models.NamespaceQuota{
		Name:           request.Name,
		Description:    request.Description,
		ProjectID:      &request.ProjectID,
		ProjectQuotaID: &request.ProjectQuotaID,
		NodeID:         request.NodeID,
	}

	if err := u.quotaRepo.CreateNamespaceQuota(namespaceQuota); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create namespace quota: %w", err))
	}

	for _, r := range request.Resources {
		resourceQuantity := models.ResourceQuantity{
			NamespaceQuotaID: &namespaceQuota.ID,
			ResourcePropID:   quotaResourcesMap[r.ResourceID].ResourceProp.ID,
			Quantity:         r.Quantity,
		}

		if err := u.quotaRepo.CreateResourceQuantity(&resourceQuantity); err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
		}
	}

	return namespaceQuota, nil
}

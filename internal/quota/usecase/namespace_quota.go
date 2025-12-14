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
			OrganizationName: org.Name,
			ProjectID:        *quota.ProjectID,
			Resources:        quota.Resources,
		})
	}
	return response, nil
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
			OrganizationName: org.Name,
			ProjectID:        *quota.ProjectID,
			Resources:        quota.Resources,
		})
	}
	return response, nil
}

func (u *QuotaUsecase) validateNamespaceQuotaRequest(request *dtos.CreateNamespaceQuotaRequest, quotaResourcesMap map[uuid.UUID]models.ResourceQuantity) error {
	if len(request.Resources) == 0 {
		return apiError.NewBadRequestError(errors.New("at least one resource quota is required"))
	}

	seenResources := make(map[uuid.UUID]struct{})
	for _, r := range request.Resources {
		if _, exists := quotaResourcesMap[r.ResourceID]; !exists {
			return apiError.NewBadRequestError(fmt.Errorf("resource %s not found in project quota", r.ResourceID))
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

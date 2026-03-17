package usecase

import (
	"errors"
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) CreateOrganizationQuota(request *dtos.CreateOrganizationQuotaRequest, userID uuid.UUID) (*models.OrganizationQuota, error) {
	if err := u.isOrgAdmin(request.FromOrganizationID, userID); err != nil {
		return nil, err
	}

	if _, err := u.orgRepo.GetOrganizationByID(request.ToOrganizationID); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to find target organization: %w", err))
	}

	org, err := u.resourceRepo.GetResourceNodeOrganization(request.NodeID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to find resource node: %w", err))
	}
	if org.ID != request.FromOrganizationID {
		return nil, apiError.NewForbiddenError(errors.New("resource node does not belong to the organization"))
	}

	if err := u.validateOrganizationQuotaRequest(request); err != nil {
		return nil, err
	}

	return u.createOrganizationQuota(request)
}

func (u *QuotaUsecase) GetOrganizationQuota(fromOrgID uuid.UUID, toOrgID uuid.UUID, userID uuid.UUID) ([]models.OrganizationQuota, error) {
	// User must be org admin of either the giving or receiving organization
	fromErr := u.isOrgAdmin(fromOrgID, userID)
	toErr := u.isOrgAdmin(toOrgID, userID)
	if fromErr != nil && toErr != nil {
		return nil, apiError.NewForbiddenError(fmt.Errorf("user is not an admin of either organization"))
	}

	quotas, err := u.quotaRepo.GetOrganizationByRelationship(fromOrgID, toOrgID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get organization quotas: %w", err))
	}
	return quotas, nil
}

func (u *QuotaUsecase) GetOrganizationQuotasByOrgID(orgID uuid.UUID, userID uuid.UUID) ([]models.OrganizationQuota, error) {
	if err := u.isOrgAdmin(orgID, userID); err != nil {
		return nil, err
	}

	quotas, err := u.quotaRepo.GetOrganizationQuotasByOrgID(orgID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get organization quotas: %w", err))
	}
	return quotas, nil
}

func (u *QuotaUsecase) validateOrganizationQuotaRequest(request *dtos.CreateOrganizationQuotaRequest) error {
	if len(request.Resources) == 0 {
		return apiError.NewBadRequestError(errors.New("at least one resource quota is required"))
	}

	isOrgQuotaExist, err := u.quotaRepo.IsOrgQuotaExist(request.FromOrganizationID, request.ToOrganizationID, request.NodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to check existing quota: %w", err))
	}

	if isOrgQuotaExist {
		return apiError.NewConflictError(errors.New("quota already exists between the organizations for this node"))
	}

	resources, err := u.resourceRepo.GetResourcesByNodeID(request.NodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get resources by node ID: %w", err))
	}

	resourcesMap := make(map[uuid.UUID]models.Resource)
	for _, resource := range resources {
		resourcesMap[resource.ID] = resource
	}

	seenResources := make(map[uuid.UUID]struct{})
	for _, r := range request.Resources {
		resource, exists := resourcesMap[r.ResourceID]
		if !exists {
			return apiError.NewNotFoundError(fmt.Errorf("resource %s not found in the specified node", r.ResourceID))
		}

		if _, duplicate := seenResources[r.ResourceID]; duplicate {
			return apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
		}
		seenResources[r.ResourceID] = struct{}{}

		if r.Quantity > resource.Quantity {
			return apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, resource.Quantity, r.ResourceID))
		}
	}

	return nil
}

func (u *QuotaUsecase) createOrganizationQuota(request *dtos.CreateOrganizationQuotaRequest) (*models.OrganizationQuota, error) {
	orgQuota := &models.OrganizationQuota{
		Name:        request.Name,
		Description: request.Description,
		NodeID:      request.NodeID,
		FromOrgID:   request.FromOrganizationID,
		ToOrgID:     request.ToOrganizationID,
	}

	if err := u.quotaRepo.CreateOrgQuota(orgQuota); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create organization quota: %w", err))
	}

	for _, r := range request.Resources {
		resourceProperty := models.ResourceProperty{
			ResourceID:  r.ResourceID,
			Price:       r.Price,
			MaxDuration: r.Duration,
		}

		if err := u.quotaRepo.CreateResourceProperty(&resourceProperty); err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource property: %w", err))
		}

		resourceQuantity := models.ResourceQuantity{
			OrganizationQuotaID: &orgQuota.ID,
			Quantity:            r.Quantity,
			ResourcePropID:      resourceProperty.ID,
		}

		if err := u.quotaRepo.CreateResourceQuantity(&resourceQuantity); err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
		}
	}

	return orgQuota, nil
}

func (u *QuotaUsecase) DeleteOrganizationQuota(quotaID uuid.UUID, userID uuid.UUID) error {
	orgQuota, err := u.quotaRepo.GetOrgQuotaByID(quotaID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("organization quota not found: %w", err))
	}

	if err := u.isOrgAdmin(orgQuota.FromOrgID, userID); err != nil {
		return err
	}

	// Cascading soft-delete: org quota -> project quotas -> namespace quotas -> resource quantities
	// 1. Get all child project quotas and cascade delete each
	projectQuotaIDs, err := u.quotaRepo.SoftDeleteProjectQuotasByOrgQuotaID(quotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to cascade delete project quotas: %w", err))
	}

	for _, pqID := range projectQuotaIDs {
		// Cascade delete namespace quotas under each project quota
		nsQuotaIDs, err := u.quotaRepo.SoftDeleteNamespaceQuotasByProjectQuotaID(pqID)
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
		if err := u.quotaRepo.SoftDeleteResourceQuantitiesByProjectQuotaID(pqID); err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota resource quantities: %w", err))
		}
	}

	// Soft-delete resource quantities for the org quota itself
	if err := u.quotaRepo.SoftDeleteResourceQuantitiesByOrgQuotaID(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete org quota resource quantities: %w", err))
	}

	// Finally, soft-delete the organization quota
	if err := u.quotaRepo.DeleteOrganizationQuota(quotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete organization quota: %w", err))
	}

	return nil
}

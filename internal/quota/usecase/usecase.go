package usecase

import (
	"errors"
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	resourcesInterfaces "github.com/ClearingHouse/internal/resources/interfaces"
	"github.com/google/uuid"
)

type QuotaUsecase struct {
	quotaRepo    interfaces.QuotaRepository
	resourceRepo resourcesInterfaces.ResourceRepository
}

func NewQuotaUsecase(quotaRepo interfaces.QuotaRepository, resourceRepo resourcesInterfaces.ResourceRepository) interfaces.QuotaUsecase {
	return &QuotaUsecase{
		quotaRepo:    quotaRepo,
		resourceRepo: resourceRepo,
	}
}

func (u *QuotaUsecase) CreateOrganizationQuotaGroup(request *dtos.CreateOrganizationQuotaRequest) (*models.OrganizationQuotaGroup, error) {
	if len(request.Resources) == 0 {
		return nil, errors.New("at least one resource quota is required")
	}

	var resources []models.Resource
	var resourceTypeCheck = map[string]bool{}
	var poolID uuid.UUID

	for i, r := range request.Resources {
		resource, err := u.resourceRepo.GetResourceByID(r.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("resource not found: %w", err)
		}

		if resource.ResourcePool.OrganizationID != request.FromOrganizationID {
			return nil, errors.New("unauthorized: resource not owned by organization")
		}

		if i == 0 {
			poolID = resource.ResourcePoolID
		} else if resource.ResourcePoolID != poolID {
			return nil, errors.New("all resources must be from the same pool")
		}

		rtName := resource.ResourceType.Name
		if resourceTypeCheck[rtName] {
			return nil, fmt.Errorf("duplicate resource type: %s", rtName)
		}
		resourceTypeCheck[rtName] = true

		resources = append(resources, *resource)
	}

	existingQuotaGroup, err := u.quotaRepo.FindExistingOrganizationQuotaGroup(request.FromOrganizationID, request.ToOrganizationID, poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing quota: %w", err)
	}

	if existingQuotaGroup != nil {
		return nil, fmt.Errorf("quota already exists between the organizations for this pool")
	}

	quotaGroup := &models.OrganizationQuotaGroup{
		Name:               request.Name,
		Description:        request.Description,
		FromOrganizationID: request.FromOrganizationID,
		ToOrganizationID:   request.ToOrganizationID,
	}

	err = u.quotaRepo.CreateOrganizationQuotaGroup(quotaGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to create quota group: %w", err)
	}

	for _, resource := range request.Resources {
		resourceProperty := models.ResourceProperty{
			ResourceID: resource.ResourceID,
			Price:      resource.Price,
		}

		err := u.quotaRepo.CreateResourceProperty(&resourceProperty)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource property: %w", err)
		}

		resourceQuantity := models.ResourceQuantity{
			OrganizationQuotaGroupID: &quotaGroup.BaseModel.ID,
			ResourcePropertyID:       resourceProperty.ID,
			Quantity:                 resource.Quantity,
		}

		err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource quantity: %w", err)
		}
	}

	return quotaGroup, nil
}

func (u *QuotaUsecase) FindOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuotaGroup, error) {
	quotaGroups, err := u.quotaRepo.FindOrganizationQuotaGroup(fromOrgID, toOrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to find organization quota groups: %w", err)
	}

	return quotaGroups, nil
}

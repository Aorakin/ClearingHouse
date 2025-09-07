package usecase

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/resources/dtos"
	"github.com/ClearingHouse/internal/resources/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceUsecase struct {
	poolRepo         interfaces.ResourcePoolRepository
	resourceRepo     interfaces.ResourceRepository
	resourceTypeRepo interfaces.ResourceTypeRepository
}

func NewResourceUsecase(poolRepo interfaces.ResourcePoolRepository, resourceRepo interfaces.ResourceRepository, resourceTypeRepo interfaces.ResourceTypeRepository) interfaces.ResourceUsecase {
	return &ResourceUsecase{
		poolRepo:         poolRepo,
		resourceRepo:     resourceRepo,
		resourceTypeRepo: resourceTypeRepo,
	}
}

func (u *ResourceUsecase) GetResources(orgID uuid.UUID) ([]dtos.ResourcePoolResponse, error) {
	resourcePools, err := u.poolRepo.GetResourcePoolByOrgID(orgID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No resource pools found for the organization
		}
		return nil, err
	}

	var response []dtos.ResourcePoolResponse
	for _, pool := range resourcePools {
		response = append(response, dtos.NewResourcePoolResponse(pool))

	}

	return response, nil
}

func (u *ResourceUsecase) CreateResourcePool(request *dtos.CreateResourcePoolRequest) (*models.ResourcePool, error) {
	resourcePool := &models.ResourcePool{
		OrganizationID: request.OrganizationID,
		Name:           request.Name,
	}
	resourcePool, err := u.poolRepo.CreateResourcePool(resourcePool)
	if err != nil {
		return nil, err
	}
	return resourcePool, nil
}

func (u *ResourceUsecase) GetResourceTypes() ([]models.ResourceType, error) {
	resourceTypes, err := u.resourceTypeRepo.GetResourceTypes()
	if err != nil {
		return nil, err
	}
	return resourceTypes, nil
}

func (u *ResourceUsecase) CreateResourceType(request *dtos.CreateResourceTypeRequest) (*models.ResourceType, error) {
	resourceType := &models.ResourceType{
		Unit: request.Unit,
		Name: request.Name,
	}
	resourceType, err := u.resourceTypeRepo.CreateResourceType(resourceType)
	if err != nil {
		return nil, err
	}
	return resourceType, nil
}

func (u *ResourceUsecase) CreateResource(request *dtos.CreateResourceRequest) (*models.Resource, error) {
	resource := &models.Resource{
		ResourcePoolID: request.ResourcePoolID,
		ResourceTypeID: request.ResourceTypeID,
		Quantity:       request.Quantity,
		Name:           request.Name,
	}
	resource, err := u.resourceRepo.CreateResource(resource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (u *ResourceUsecase) UpdateResource(resourceID uuid.UUID, request *dtos.UpdateResourceRequest) (*models.Resource, error) {
	resource, err := u.resourceRepo.GetResourceByID(resourceID)
	if err != nil {
		return nil, err
	}

	resource.Quantity = request.Quantity
	resource.Name = request.Name

	updatedResource, err := u.resourceRepo.UpdateResource(resource)
	if err != nil {
		return nil, err
	}
	return updatedResource, nil
}

func (u *ResourceUsecase) GetResourceProperty(resourceID uuid.UUID) (*models.Resource, error) {
	resource, err := u.resourceRepo.GetResourceByID(resourceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return resource, nil
}

func (u *ResourceUsecase) GetResourcePool(resourcePoolID *uuid.UUID) (*models.ResourcePool, error) {
	resourcePool, err := u.poolRepo.GetResourcePoolByID(*resourcePoolID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return resourcePool, nil
}

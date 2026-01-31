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

func (u *ResourceUsecase) CreateResourcePool(request *dtos.CreateResourcePoolRequest) (*models.ResourcePool, error) {
	resourcePool := &models.ResourcePool{
		OrganizationID: request.OrganizationID,
		Name:           request.Name,
		GlideletURN:    request.GlideletURN,
	}
	resourcePool, err := u.poolRepo.CreateResourcePool(resourcePool)
	if err != nil {
		return nil, err
	}
	return resourcePool, nil
}

func (u *ResourceUsecase) CreateResourceNode(request *dtos.CreateResourceNodeRequest) (*models.ResourceNode, error) {
	resourceNode := &models.ResourceNode{
		ResourcePoolID: request.ResourcePoolID,
		Name:           request.Name,
	}
	createdNode, err := u.resourceRepo.CreateResourceNode(resourceNode)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return createdNode, nil
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
		NodeID:         request.ResourceNodeID,
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

func (u *ResourceUsecase) GetResourcePool(resourcePoolID uuid.UUID) (*models.ResourcePool, error) {
	resourcePool, err := u.poolRepo.GetResourcePoolByID(resourcePoolID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return resourcePool, nil
}

func (u *ResourceUsecase) GetResourceNode(nodeID uuid.UUID) (*models.ResourceNode, error) {
	resourceNode, err := u.resourceRepo.GetResourceNodeByID(nodeID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return resourceNode, nil
}

func (u *ResourceUsecase) GetResources(orgID uuid.UUID) ([]models.ResourcePool, error) {
	resourcePools, err := u.poolRepo.GetResourcePoolByOrgID(orgID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No resource pools found for the organization
		}
		return nil, err
	}

	return resourcePools, nil
}

func (u *ResourceUsecase) GetResourceTypes() ([]models.ResourceType, error) {
	resourceTypes, err := u.resourceTypeRepo.GetResourceTypes()
	if err != nil {
		return nil, err
	}
	return resourceTypes, nil
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

func (u *ResourceUsecase) DeleteResourcePool(resourcePoolID uuid.UUID) error {
	// Get resource pool with all nodes
	resourcePool, err := u.poolRepo.GetResourcePoolByID(resourcePoolID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apiError.NewNotFoundError("resource pool not found")
		}
		return apiError.NewInternalServerError(err)
	}

	// Check if resource pool has any nodes
	if len(resourcePool.Nodes) > 0 {
		return apiError.NewBadRequestError("cannot delete resource pool with existing nodes. Please delete all nodes first")
	}

	// Check if resource pool is being used by any active tickets
	hasActiveTickets, err := u.poolRepo.HasActiveTickets(resourcePoolID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	if hasActiveTickets {
		return apiError.NewBadRequestError("cannot delete resource pool with active tickets. Please wait for all tickets to complete")
	}

	// Soft delete the resource pool
	if err := u.poolRepo.DeleteResourcePool(resourcePoolID); err != nil {
		return apiError.NewInternalServerError(err)
	}

	return nil
}

func (u *ResourceUsecase) DeleteResourceNode(nodeID uuid.UUID) error {
	// Get resource node with all resources
	resourceNode, err := u.resourceRepo.GetResourceNodeByID(nodeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apiError.NewNotFoundError("resource node not found")
		}
		return apiError.NewInternalServerError(err)
	}

	// Check if node has any resources
	if len(resourceNode.Resources) > 0 {
		return apiError.NewBadRequestError("cannot delete resource node with existing resources. Please delete all resources first")
	}

	// Check if node is being used by any active tickets
	hasActiveTickets, err := u.resourceRepo.HasActiveTicketsByNodeID(nodeID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	if hasActiveTickets {
		return apiError.NewBadRequestError("cannot delete resource node with active tickets. Please wait for all tickets to complete")
	}

	// Soft delete the resource node
	if err := u.resourceRepo.DeleteResourceNode(nodeID); err != nil {
		return apiError.NewInternalServerError(err)
	}

	return nil
}

func (u *ResourceUsecase) DeleteResource(resourceID uuid.UUID) error {
	// Get resource to check if it exists
	resource, err := u.resourceRepo.GetResourceByID(resourceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apiError.NewNotFoundError("resource not found")
		}
		return apiError.NewInternalServerError(err)
	}

	// Check if resource has any resource properties
	hasResourceProperties, err := u.resourceRepo.HasResourceProperties(resourceID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	if hasResourceProperties {
		return apiError.NewBadRequestError("cannot delete resource with existing resource properties and quotas. Please remove all quota allocations first")
	}

	// Check if resource is being used by any active tickets
	hasActiveTickets, err := u.resourceRepo.HasActiveTicketsByResourceID(resourceID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	if hasActiveTickets {
		return apiError.NewBadRequestError("cannot delete resource with active tickets. Please wait for all tickets to complete")
	}

	// Soft delete the resource
	if err := u.resourceRepo.DeleteResource(resource.ID); err != nil {
		return apiError.NewInternalServerError(err)
	}

	return nil
}

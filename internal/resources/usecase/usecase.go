package usecase

import (
	"fmt"

	"github.com/ClearingHouse/internal/models"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
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
	quotaRepo        quotaInterfaces.QuotaRepository
}

func NewResourceUsecase(poolRepo interfaces.ResourcePoolRepository, resourceRepo interfaces.ResourceRepository, resourceTypeRepo interfaces.ResourceTypeRepository, quotaRepo quotaInterfaces.QuotaRepository) interfaces.ResourceUsecase {
	return &ResourceUsecase{
		poolRepo:         poolRepo,
		resourceRepo:     resourceRepo,
		resourceTypeRepo: resourceTypeRepo,
		quotaRepo:        quotaRepo,
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
		DisplayName:    request.DisplayName,
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

func (u *ResourceUsecase) UpdateResourcePool(resourcePoolID uuid.UUID, request *dtos.UpdateResourcePoolRequest) (*models.ResourcePool, error) {
	// Get existing resource pool
	resourcePool, err := u.poolRepo.GetResourcePoolByID(resourcePoolID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apiError.NewNotFoundError("resource pool not found")
		}
		return nil, apiError.NewInternalServerError(err)
	}

	// Update fields
	resourcePool.Name = request.Name
	resourcePool.GlideletURN = request.GlideletURN

	// Save changes
	updatedPool, err := u.poolRepo.UpdateResourcePool(resourcePool)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return updatedPool, nil
}

func (u *ResourceUsecase) GetResourceNode(nodeID uuid.UUID) (*models.ResourceNode, error) {
	resourceNode, err := u.resourceRepo.GetResourceNodeByID(nodeID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return resourceNode, nil
}

func (u *ResourceUsecase) UpdateResourceNode(nodeID uuid.UUID, request *dtos.UpdateResourceNodeRequest) (*models.ResourceNode, error) {
	// Get existing resource node
	resourceNode, err := u.resourceRepo.GetResourceNodeByID(nodeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apiError.NewNotFoundError("resource node not found")
		}
		return nil, apiError.NewInternalServerError(err)
	}

	// Update fields
	resourceNode.Name = request.Name
	resourceNode.DisplayName = request.DisplayName

	// Save changes
	updatedNode, err := u.resourceRepo.UpdateResourceNode(resourceNode)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return updatedNode, nil
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
	oldQuantity := resource.Quantity

	resource.Quantity = request.Quantity
	resource.Name = request.Name

	updatedResource, err := u.resourceRepo.UpdateResource(resource)
	if err != nil {
		return nil, err
	}

	// Clamp downstream quotas when a resource capacity is reduced.
	if request.Quantity < oldQuantity {
		if err := u.cascadeResourceCapacityUpdate(resource.ID, resource.NodeID, request.Quantity); err != nil {
			return nil, err
		}
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
	// resourceNode, err := u.resourceRepo.GetResourceNodeByID(nodeID)
	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		return apiError.NewNotFoundError("resource node not found")
	// 	}
	// 	return apiError.NewInternalServerError(err)
	// }

	// Check if node has any resources
	// if len(resourceNode.Resources) > 0 {
	// 	return apiError.NewBadRequestError("cannot delete resource node with existing resources. Please delete all resources first")
	// }

	// Check if node is being used by any active tickets
	hasActiveTickets, err := u.resourceRepo.HasActiveTicketsByNodeID(nodeID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	if hasActiveTickets {
		return apiError.NewBadRequestError("cannot delete resource node with active tickets. Please wait for all tickets to complete")
	}

	if err := u.cascadeDeleteNodeQuotas(nodeID); err != nil {
		return err
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

func (u *ResourceUsecase) cascadeResourceCapacityUpdate(resourceID uuid.UUID, nodeID uuid.UUID, newLimit uint) error {
	node, err := u.resourceRepo.GetResourceNodeByID(nodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get resource node for quota cascade: %w", err))
	}

	orgQuotas, err := u.quotaRepo.GetOrganizationQuotasByFromOrgAndNode(node.ResourcePool.OrganizationID, nodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get organization quotas for cascade: %w", err))
	}

	for _, orgQuota := range orgQuotas {
		if err := u.clampOrgQuotaAndChildren(orgQuota.ID, resourceID, newLimit); err != nil {
			return err
		}
	}

	internalProjectQuotas, err := u.quotaRepo.GetInternalProjectQuotasByOrgAndNode(node.ResourcePool.OrganizationID, nodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get internal project quotas for cascade: %w", err))
	}

	for _, projectQuota := range internalProjectQuotas {
		if err := u.clampProjectQuotaAndNamespaces(projectQuota.ID, resourceID, newLimit); err != nil {
			return err
		}
	}

	return nil
}

func (u *ResourceUsecase) clampOrgQuotaAndChildren(orgQuotaID uuid.UUID, resourceID uuid.UUID, newLimit uint) error {
	orgResources, err := u.quotaRepo.GetResourceQuantitiesByOrgQuotaID(orgQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get organization quota resources: %w", err))
	}

	for _, rq := range orgResources {
		if rq.ResourceProp.ResourceID != resourceID {
			continue
		}

		if rq.Quantity > newLimit {
			if err := u.quotaRepo.UpdateResourceQuantity(rq.ID, newLimit); err != nil {
				return apiError.NewInternalServerError(fmt.Errorf("failed to clamp organization quota resource: %w", err))
			}
		}

		if err := u.cascadeOrgQuotaToProjects(orgQuotaID, resourceID, newLimit); err != nil {
			return err
		}

		break
	}

	return nil
}

func (u *ResourceUsecase) cascadeOrgQuotaToProjects(orgQuotaID uuid.UUID, resourceID uuid.UUID, newLimit uint) error {
	projectQuotas, err := u.quotaRepo.GetProjectQuotasByOrgQuotaID(orgQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get project quotas from organization quota: %w", err))
	}

	for _, projectQuota := range projectQuotas {
		if err := u.clampProjectQuotaAndNamespaces(projectQuota.ID, resourceID, newLimit); err != nil {
			return err
		}
	}

	return nil
}

func (u *ResourceUsecase) clampProjectQuotaAndNamespaces(projectQuotaID uuid.UUID, resourceID uuid.UUID, newLimit uint) error {
	projectResources, err := u.quotaRepo.GetResourceQuantitiesByProjectQuotaID(projectQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get project quota resources: %w", err))
	}

	for _, rq := range projectResources {
		if rq.ResourceProp.ResourceID != resourceID {
			continue
		}

		if rq.Quantity > newLimit {
			if err := u.quotaRepo.UpdateResourceQuantity(rq.ID, newLimit); err != nil {
				return apiError.NewInternalServerError(fmt.Errorf("failed to clamp project quota resource: %w", err))
			}
		}

		if err := u.cascadeProjectQuotaToNamespaces(projectQuotaID, resourceID, newLimit); err != nil {
			return err
		}

		break
	}

	return nil
}

func (u *ResourceUsecase) cascadeProjectQuotaToNamespaces(projectQuotaID uuid.UUID, resourceID uuid.UUID, newLimit uint) error {
	namespaceQuotas, err := u.quotaRepo.GetNamespaceQuotasByProjectQuotaID(projectQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quotas from project quota: %w", err))
	}

	for _, namespaceQuota := range namespaceQuotas {
		for _, rq := range namespaceQuota.Resources {
			if rq.ResourceProp.ResourceID == resourceID && rq.Quantity > newLimit {
				if err := u.quotaRepo.UpdateResourceQuantity(rq.ID, newLimit); err != nil {
					return apiError.NewInternalServerError(fmt.Errorf("failed to clamp namespace quota resource: %w", err))
				}
			}
		}
	}

	return nil
}

func (u *ResourceUsecase) cascadeDeleteNodeQuotas(nodeID uuid.UUID) error {
	node, err := u.resourceRepo.GetResourceNodeByID(nodeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apiError.NewNotFoundError("resource node not found")
		}
		return apiError.NewInternalServerError(fmt.Errorf("failed to get resource node for quota cleanup: %w", err))
	}

	orgQuotas, err := u.quotaRepo.GetOrganizationQuotasByFromOrgAndNode(node.ResourcePool.OrganizationID, nodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to list organization quotas for node cleanup: %w", err))
	}

	for _, orgQuota := range orgQuotas {
		if err := u.deleteOrganizationQuotaCascade(orgQuota.ID); err != nil {
			return err
		}
	}

	internalProjectQuotas, err := u.quotaRepo.GetInternalProjectQuotasByOrgAndNode(node.ResourcePool.OrganizationID, nodeID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to list internal project quotas for node cleanup: %w", err))
	}

	for _, projectQuota := range internalProjectQuotas {
		if err := u.deleteProjectQuotaCascade(projectQuota.ID); err != nil {
			return err
		}
	}

	return nil
}

func (u *ResourceUsecase) deleteOrganizationQuotaCascade(orgQuotaID uuid.UUID) error {
	if err := u.quotaRepo.SoftDeleteResourceQuantitiesByOrgQuotaID(orgQuotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete organization quota resources: %w", err))
	}

	projectQuotaIDs, err := u.quotaRepo.SoftDeleteProjectQuotasByOrgQuotaID(orgQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quotas from organization quota: %w", err))
	}

	for _, projectQuotaID := range projectQuotaIDs {
		namespaceQuotaIDs, err := u.quotaRepo.SoftDeleteNamespaceQuotasByProjectQuotaID(projectQuotaID)
		if err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to delete namespace quotas from project quota: %w", err))
		}

		for _, namespaceQuotaID := range namespaceQuotaIDs {
			if err := u.quotaRepo.DeleteResourceQuantitiesByNamespaceQuotaID(namespaceQuotaID); err != nil {
				return apiError.NewInternalServerError(fmt.Errorf("failed to delete namespace quota resources: %w", err))
			}
		}

		if err := u.quotaRepo.UnassignQuotaTemplatesByNamespaceQuotaIDs(namespaceQuotaIDs); err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to unassign quota templates from deleted namespace quotas: %w", err))
		}

		if err := u.quotaRepo.SoftDeleteResourceQuantitiesByProjectQuotaID(projectQuotaID); err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota resources: %w", err))
		}
	}

	if err := u.quotaRepo.DeleteOrganizationQuota(orgQuotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete organization quota: %w", err))
	}

	return nil
}

func (u *ResourceUsecase) deleteProjectQuotaCascade(projectQuotaID uuid.UUID) error {
	namespaceQuotaIDs, err := u.quotaRepo.SoftDeleteNamespaceQuotasByProjectQuotaID(projectQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete namespace quotas from project quota: %w", err))
	}

	for _, namespaceQuotaID := range namespaceQuotaIDs {
		if err := u.quotaRepo.DeleteResourceQuantitiesByNamespaceQuotaID(namespaceQuotaID); err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to delete namespace quota resources: %w", err))
		}
	}

	if err := u.quotaRepo.UnassignQuotaTemplatesByNamespaceQuotaIDs(namespaceQuotaIDs); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to unassign quota templates from deleted namespace quotas: %w", err))
	}

	if err := u.quotaRepo.SoftDeleteResourceQuantitiesByProjectQuotaID(projectQuotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota resources: %w", err))
	}

	if err := u.quotaRepo.DeleteProjectQuota(projectQuotaID); err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to delete project quota: %w", err))
	}

	return nil
}

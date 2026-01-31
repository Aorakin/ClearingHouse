package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type ResourceRepository interface {
	CreateResource(resource *models.Resource) (*models.Resource, error)
	UpdateResource(resource *models.Resource) (*models.Resource, error)
	GetResourceByID(id uuid.UUID) (*models.Resource, error)
	GetResourcePoolByID(id uuid.UUID) (*models.ResourcePool, error)
	GetResourcesByOrganizationID(orgID uuid.UUID) ([]models.Resource, error)
	GetResourcesByNodeID(nodeID uuid.UUID) ([]models.Resource, error)

	CreateResourceNode(resourceNode *models.ResourceNode) (*models.ResourceNode, error)
	GetResourceNodeByID(nodeID uuid.UUID) (*models.ResourceNode, error)
	GetResourceNodeOrganization(nodeID uuid.UUID) (*models.Organization, error)
}

type ResourceTypeRepository interface {
	GetResourceTypes() ([]models.ResourceType, error)
	CreateResourceType(resourceType *models.ResourceType) (*models.ResourceType, error)
}

type ResourcePoolRepository interface {
	GetResourcePoolByID(id uuid.UUID) (*models.ResourcePool, error)
	GetResourcePoolByOrgID(orgID uuid.UUID) ([]models.ResourcePool, error)
	CreateResourcePool(resourcePool *models.ResourcePool) (*models.ResourcePool, error)
	DeleteResourcePool(id uuid.UUID) error
	HasActiveTickets(resourcePoolID uuid.UUID) (bool, error)
}

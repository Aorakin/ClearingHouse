package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type ResourceRepository interface {
	CreateResource(resource *models.Resource) (*models.Resource, error)
	UpdateResource(resource *models.Resource) (*models.Resource, error)
	GetResourceByID(id uuid.UUID) (*models.Resource, error)
}

type ResourceTypeRepository interface {
	GetResourceTypes() ([]models.ResourceType, error)
	CreateResourceType(resourceType *models.ResourceType) (*models.ResourceType, error)
}

type ResourcePoolRepository interface {
	GetResourcePoolByOrgID(orgID uuid.UUID) ([]models.ResourcePool, error)
	CreateResourcePool(resourcePool *models.ResourcePool) (*models.ResourcePool, error)
}

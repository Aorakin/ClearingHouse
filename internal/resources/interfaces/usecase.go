package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/resources/dtos"
	"github.com/google/uuid"
)

type ResourceUsecase interface {
	GetResources(orgID uuid.UUID) ([]models.ResourcePool, error)
	GetResourceTypes() ([]models.ResourceType, error)
	CreateResource(request *dtos.CreateResourceRequest) (*models.Resource, error)
	CreateResourceType(request *dtos.CreateResourceTypeRequest) (*models.ResourceType, error)
	CreateResourcePool(request *dtos.CreateResourcePoolRequest) (*models.ResourcePool, error)
	UpdateResource(resourceID uuid.UUID, request *dtos.UpdateResourceRequest) (*models.Resource, error)
	GetResourceProperty(resourceID uuid.UUID) (*models.Resource, error)
	DeleteResource(resourceID uuid.UUID) error
	GetResourcePool(resourcePoolID uuid.UUID) (*models.ResourcePool, error)
	UpdateResourcePool(resourcePoolID uuid.UUID, request *dtos.UpdateResourcePoolRequest) (*models.ResourcePool, error)
	DeleteResourcePool(resourcePoolID uuid.UUID) error
	GetResourceNode(nodeID uuid.UUID) (*models.ResourceNode, error)
	CreateResourceNode(request *dtos.CreateResourceNodeRequest) (*models.ResourceNode, error)
	UpdateResourceNode(nodeID uuid.UUID, request *dtos.UpdateResourceNodeRequest) (*models.ResourceNode, error)
	DeleteResourceNode(nodeID uuid.UUID) error
}

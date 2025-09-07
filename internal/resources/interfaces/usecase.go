package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/resources/dtos"
	"github.com/google/uuid"
)

type ResourceUsecase interface {
	GetResources(orgID uuid.UUID) ([]dtos.ResourcePoolResponse, error)
	GetResourceTypes() ([]models.ResourceType, error)
	CreateResource(request *dtos.CreateResourceRequest) (*models.Resource, error)
	CreateResourceType(request *dtos.CreateResourceTypeRequest) (*models.ResourceType, error)
	CreateResourcePool(request *dtos.CreateResourcePoolRequest) (*models.ResourcePool, error)
	UpdateResource(resourceID uuid.UUID, request *dtos.UpdateResourceRequest) (*models.Resource, error)
	GetResourceProperty(resourceID uuid.UUID) (*models.Resource, error)
	GetResourcePool(resourcePoolID *uuid.UUID) (*models.ResourcePool, error)
}

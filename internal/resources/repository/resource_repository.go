package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceRepository struct {
	db *gorm.DB
}

func (r *ResourceRepository) CreateResource(resource *models.Resource) (*models.Resource, error) {
	if err := r.db.Create(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func (r *ResourceRepository) UpdateResource(resource *models.Resource) (*models.Resource, error) {
	if err := r.db.Save(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func (r *ResourceRepository) GetResourceByID(id uuid.UUID) (*models.Resource, error) {
	var resource models.Resource
	if err := r.db.Preload("ResourceType").Preload("ResourcePool").First(&resource, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *ResourceRepository) GetResourcePoolByID(id uuid.UUID) (*models.ResourcePool, error) {
	var resourcePool models.ResourcePool
	err := r.db.Preload("Resources.ResourceType").First(&resourcePool, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &resourcePool, nil
}

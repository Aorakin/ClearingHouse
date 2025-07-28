package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourcePoolRepository struct {
	db *gorm.DB
}

func (r *ResourcePoolRepository) GetResourcePoolByOrgID(orgID uuid.UUID) ([]models.ResourcePool, error) {
	var resourcePools []models.ResourcePool
	err := r.db.
		Preload("Resources").
		Preload("Resources.ResourceType").
		Where("organization_id = ?", orgID).
		Find(&resourcePools).Error

	if err != nil {
		return nil, err
	}
	return resourcePools, nil
}

func (r *ResourcePoolRepository) CreateResourcePool(resourcePool *models.ResourcePool) (*models.ResourcePool, error) {
	if err := r.db.Create(resourcePool).Error; err != nil {
		return nil, err
	}
	return resourcePool, nil
}

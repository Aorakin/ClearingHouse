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
		Preload("Nodes.Resources.ResourceType").
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

func (r *ResourcePoolRepository) GetResourcePoolByID(id uuid.UUID) (*models.ResourcePool, error) {
	var resourcePool models.ResourcePool
	err := r.db.Preload("Nodes.Resources.ResourceType").First(&resourcePool, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &resourcePool, nil
}

func (r *ResourcePoolRepository) DeleteResourcePool(id uuid.UUID) error {
	// GORM will perform soft delete automatically (sets deleted_at)
	return r.db.Delete(&models.ResourcePool{}, "id = ?", id).Error
}

func (r *ResourcePoolRepository) HasActiveTickets(resourcePoolID uuid.UUID) (bool, error) {
	var count int64

	// Check for tickets with active statuses: pending, using, redeeming
	err := r.db.Table("tickets").
		Where("resource_pool_id = ? AND status IN ?", resourcePoolID, []string{"pending", "using", "redeeming"}).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

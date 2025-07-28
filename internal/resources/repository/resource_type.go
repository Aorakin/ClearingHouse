package repository

import (
	"github.com/ClearingHouse/internal/models"
	"gorm.io/gorm"
)

type ResourceTypeRepository struct {
	db *gorm.DB
}

func (r *ResourceTypeRepository) CreateResourceType(resourceType *models.ResourceType) (*models.ResourceType, error) {
	if err := r.db.Create(resourceType).Error; err != nil {
		return nil, err
	}
	return resourceType, nil
}

func (r *ResourceTypeRepository) GetResourceTypes() ([]models.ResourceType, error) {
	var resourceTypes []models.ResourceType
	if err := r.db.Find(&resourceTypes).Error; err != nil {
		return nil, err
	}
	return resourceTypes, nil
}

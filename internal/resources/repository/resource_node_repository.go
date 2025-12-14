package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceNodeRepository struct {
	db *gorm.DB
}

func (r *ResourceNodeRepository) CreateResourceNode(resourceNode *models.ResourceNode) (*models.ResourceNode, error) {
	if err := r.db.Create(resourceNode).Error; err != nil {
		return nil, err
	}
	return resourceNode, nil
}

func (r *ResourceNodeRepository) GetResourceNodeByID(nodeID uuid.UUID) (*models.ResourceNode, error) {
	var resourceNode models.ResourceNode
	if err := r.db.First(&resourceNode, "id = ?", nodeID).Error; err != nil {
		return nil, err
	}
	return &resourceNode, nil
}

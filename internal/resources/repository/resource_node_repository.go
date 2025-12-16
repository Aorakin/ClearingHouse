package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func (r *ResourceRepository) CreateResourceNode(resourceNode *models.ResourceNode) (*models.ResourceNode, error) {
	if err := r.db.Create(resourceNode).Error; err != nil {
		return nil, err
	}
	return resourceNode, nil
}

func (r *ResourceRepository) GetResourceNodeByID(nodeID uuid.UUID) (*models.ResourceNode, error) {
	var resourceNode models.ResourceNode
	if err := r.db.Preload("Resources").Preload("ResourcePool").First(&resourceNode, "id = ?", nodeID).Error; err != nil {
		return nil, err
	}
	return &resourceNode, nil
}

func (r *ResourceRepository) GetResourceNodeOrganization(nodeID uuid.UUID) (*models.Organization, error) {
	var resourceNode models.ResourceNode
	if err := r.db.Preload("ResourcePool.Organization").First(&resourceNode, "id = ?", nodeID).Error; err != nil {
		return nil, err
	}
	return &resourceNode.ResourcePool.Organization, nil
}

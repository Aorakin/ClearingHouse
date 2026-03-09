package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ResourceRepository) CreateResourceNode(resourceNode *models.ResourceNode) (*models.ResourceNode, error) {
	if err := r.db.Create(resourceNode).Error; err != nil {
		return nil, err
	}
	return resourceNode, nil
}

func (r *ResourceRepository) GetResourceNodeByID(nodeID uuid.UUID) (*models.ResourceNode, error) {
	var resourceNode models.ResourceNode
	if err := r.db.
		Preload("Resources", func(db *gorm.DB) *gorm.DB {
			return db.Order("name")
		}).
		Preload("Resources.ResourceType").
		Preload("ResourcePool").
		First(&resourceNode, "id = ?", nodeID).Error; err != nil {
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

func (r *ResourceRepository) UpdateResourceNode(resourceNode *models.ResourceNode) (*models.ResourceNode, error) {
	if err := r.db.Save(resourceNode).Error; err != nil {
		return nil, err
	}
	return resourceNode, nil
}

func (r *ResourceRepository) DeleteResourceNode(nodeID uuid.UUID) error {
	// GORM will perform soft delete automatically (sets deleted_at)
	return r.db.Delete(&models.ResourceNode{}, "id = ?", nodeID).Error
}

func (r *ResourceRepository) HasActiveTicketsByNodeID(nodeID uuid.UUID) (bool, error) {
	var count int64

	// Check for tickets with active statuses: pending, using, redeeming
	err := r.db.Table("tickets").
		Where("node_id = ? AND status IN ?", nodeID, []string{"pending", "using", "redeeming"}).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

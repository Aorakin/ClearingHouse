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
	if err := r.db.Preload("ResourceType").First(&resource, "id = ?", id).Error; err != nil {
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

func (r *ResourceRepository) GetResourcesByOrganizationID(orgID uuid.UUID) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Joins("JOIN resource_pools ON resource_pools.id = resources.resource_pool_id").
		Where("resource_pools.organization_id = ?", orgID).
		Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *ResourceRepository) GetResourcesByNodeID(nodeID uuid.UUID) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Where("node_id = ?", nodeID).Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *ResourceRepository) DeleteResource(resourceID uuid.UUID) error {
	// GORM will perform soft delete automatically (sets deleted_at)
	return r.db.Delete(&models.Resource{}, "id = ?", resourceID).Error
}

func (r *ResourceRepository) HasResourceProperties(resourceID uuid.UUID) (bool, error) {
	var count int64

	// Check if resource has any resource properties with quantities that are still in active quotas
	// Need to check all three types of quotas: organization, project, and namespace
	// Also need to check that related projects are not deleted
	err := r.db.Table("resource_properties rp").
		Joins("JOIN resource_quantities rq ON rq.resource_prop_id = rp.id").
		Joins("LEFT JOIN organization_quota oq ON oq.id = rq.organization_quota_id").
		Joins("LEFT JOIN project_quota pq ON pq.id = rq.project_quota_id").
		Joins("LEFT JOIN projects p ON p.id = pq.project_id").
		Joins("LEFT JOIN namespace_quota nq ON nq.id = rq.namespace_quota_id").
		Where("rp.resource_id = ?", resourceID).
		Where("(rq.organization_quota_id IS NOT NULL AND oq.deleted_at IS NULL) OR " +
			"(rq.project_quota_id IS NOT NULL AND pq.deleted_at IS NULL AND p.deleted_at IS NULL) OR " +
			"(rq.namespace_quota_id IS NOT NULL AND nq.deleted_at IS NULL)").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ResourceRepository) HasActiveTicketsByResourceID(resourceID uuid.UUID) (bool, error) {
	var count int64

	// Check for tickets with active statuses through ticket_resources
	err := r.db.Table("ticket_resources tr").
		Joins("JOIN tickets t ON t.id = tr.ticket_id").
		Where("tr.resource_id = ? AND t.status IN ?", resourceID, []string{"pending", "using", "redeeming"}).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

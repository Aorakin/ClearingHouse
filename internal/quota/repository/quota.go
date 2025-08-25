package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuotaRepository struct {
	db *gorm.DB
}

func NewQuotaRepository(db *gorm.DB) interfaces.QuotaRepository {
	return &QuotaRepository{db: db}
}
func (r *QuotaRepository) FindOrganizationQuotaGroup(fromOrgId uuid.UUID, toOrgId uuid.UUID) ([]models.OrganizationQuotaGroup, error) {
	var quotaGroups []models.OrganizationQuotaGroup
	err := r.db.Debug().Preload("Resources.ResourceProperty").Where("from_organization_id = ? AND to_organization_id = ?", fromOrgId, toOrgId).Find(&quotaGroups).Error
	if err != nil {
		return nil, err
	}
	return quotaGroups, nil
}

func (r *QuotaRepository) FindExistingOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID, poolID uuid.UUID) (*models.OrganizationQuotaGroup, error) {
	var quotaGroup models.OrganizationQuotaGroup
	err := r.db.Joins("JOIN resource_quantities rq ON rq.organization_quota_group_id = organization_quota_groups.id").
		Joins("JOIN resource_properties rp ON rq.resource_property_id = rp.id").
		Joins("JOIN resources r ON rp.resource_id = r.id").
		Where("from_organization_id = ? AND to_organization_id = ? AND r.resource_pool_id = ?", fromOrgID, toOrgID, poolID).
		First(&quotaGroup).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No existing quota group found
		}
		return nil, err // Other error
	}
	return &quotaGroup, nil
}

func (r *QuotaRepository) CreateOrganizationQuotaGroup(quota *models.OrganizationQuotaGroup) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error {
	return r.db.Create(resourceQuantity).Error
}

func (r *QuotaRepository) CreateResourceProperty(resourceProperty *models.ResourceProperty) error {
	return r.db.Create(resourceProperty).Error
}

func (r *QuotaRepository) CreateProjectQuotaGroup(quota *models.ProjectQuotaGroup) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) CreateNamespaceQuotaGroup(quota *models.NamespaceQuotaGroup) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) FindOrganizationQuotaGroupByID(id uuid.UUID) (*models.OrganizationQuotaGroup, error) {
	var quotaGroup models.OrganizationQuotaGroup
	err := r.db.Preload("Resources.ResourceProperty").Where("id = ?", id).Find(&quotaGroup).Error
	if err != nil {
		return nil, err
	}
	return &quotaGroup, nil
}

func (r *QuotaRepository) GetOrgUsage(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error) {
	var total uint

	err := r.db.Debug().
		Table("resource_quantities rq").
		Joins("JOIN resource_properties rp ON rp.id = rq.resource_property_id").
		Joins("LEFT JOIN project_quota_groups pqg ON pqg.id = rq.project_quota_group_id").
		Where("rp.resource_id = ? AND pqg.organization_quota_group_id = ?", resourceID, orgQuotaGroupID).
		Select("COALESCE(SUM(rq.quantity), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *QuotaRepository) GetOrgQuotaQuantity(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error) {
	var total uint
	err := r.db.
		Model(&models.ResourceQuantity{}).
		Joins("JOIN resource_properties rp ON rp.id = resource_quantities.resource_property_id").
		Where("organization_quota_group_id = ? AND rp.resource_id = ?", orgQuotaGroupID, resourceID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *QuotaRepository) GetProjectQuotaQuantity(projQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error) {
	var total uint
	err := r.db.
		Model(&models.ResourceQuantity{}).
		Joins("JOIN resource_properties rp ON rp.id = resource_quantities.resource_property_id").
		Where("project_quota_group_id = ? AND rp.resource_id = ?", projQuotaGroupID, resourceID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *QuotaRepository) GetResourcePropertyByOrg(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (*models.ResourceProperty, error) {
	var resourceProperty models.ResourceProperty
	err := r.db.
		Model(&models.ResourceProperty{}).
		Joins(`
			JOIN resource_quantities rq 
				ON rq.resource_property_id = resource_properties.id
		`).
		Where("rq.organization_quota_group_id = ? AND resource_properties.resource_id = ?", orgQuotaGroupID, resourceID).
		First(&resourceProperty).Error

	if err != nil {
		return nil, err
	}
	return &resourceProperty, nil
}

func (r *QuotaRepository) GetResourcePropertyByProj(projQuotaGroupID uuid.UUID, resourceID uuid.UUID) (*models.ResourceProperty, error) {
	var resourceProperty models.ResourceProperty
	err := r.db.
		Model(&models.ResourceProperty{}).
		Joins(`
			JOIN resource_quantities rq 
				ON rq.resource_property_id = resource_properties.id
		`).
		Where("rq.project_quota_group_id = ? AND resource_properties.resource_id = ?", projQuotaGroupID, resourceID).
		First(&resourceProperty).Error

	if err != nil {
		return nil, err
	}
	return &resourceProperty, nil
}

func (r *QuotaRepository) FindProjectQuotaGroupByID(id uuid.UUID) (*models.ProjectQuotaGroup, error) {
	var projectQuotaGroup models.ProjectQuotaGroup
	err := r.db.Preload("Resources.ResourceProperty").Where("id = ?", id).Find(&projectQuotaGroup).Error
	if err != nil {
		return nil, err
	}
	return &projectQuotaGroup, nil
}

func (r *QuotaRepository) FindProjectQuotaGroupByProjectID(projectID uuid.UUID) ([]models.ProjectQuotaGroup, error) {
	var projectQuotaGroups []models.ProjectQuotaGroup
	err := r.db.Joins("JOIN project_quotas pq ON pq.project_quota_group_id = project_quota_groups.id").
		Preload("Resources.ResourceProperty").Where("pq.project_id = ?", projectID).
		Find(&projectQuotaGroups).Error
	if err != nil {
		return nil, err
	}
	return projectQuotaGroups, nil
}

func (r *QuotaRepository) FindNamespaceQuotaGroupByID(id uuid.UUID) (*models.NamespaceQuotaGroup, error) {
	var namespaceQuotaGroup models.NamespaceQuotaGroup
	err := r.db.Preload("Resources.ResourceProperty").Where("id = ?", id).Find(&namespaceQuotaGroup).Error
	if err != nil {
		return nil, err
	}
	return &namespaceQuotaGroup, nil
}

func (r *QuotaRepository) AssignQuotaToNamespace(namespaceID uuid.UUID, quotaGroupID uuid.UUID) error {
	return r.db.Model(&models.Namespace{}).Where("id = ?", namespaceID).Update("quota_group_id", quotaGroupID).Error
}

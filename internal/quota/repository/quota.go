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

func (r *QuotaRepository) CreateOrgQuota(quota *models.OrganizationQuota) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) GetOrganizationByRelationship(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuota, error) {
	var organizations []models.OrganizationQuota
	err := r.db.Preload("Resources.ResourceProp").Where("from_org_id = ? AND to_org_id = ?", fromOrgID, toOrgID).
		Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *QuotaRepository) GetOrgQuotaByID(id uuid.UUID) (*models.OrganizationQuota, error) {
	var orgQuota models.OrganizationQuota
	err := r.db.Preload("Resources.ResourceProp").First(&orgQuota, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &orgQuota, nil
}

func (r *QuotaRepository) GetOrgUsage(quotaID uuid.UUID, resourceID uuid.UUID) (uint, error) {
	var total uint

	err := r.db.Debug().
		Table("resource_quantities rq").
		Joins("JOIN resource_properties rp ON rp.id = rq.resource_prop_id").
		Joins("JOIN project_quota pq ON pq.id = rq.project_quota_id").
		Where("rp.resource_id = ? AND pq.organization_quota_id = ?", resourceID, quotaID).
		Select("COALESCE(SUM(rq.quantity), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *QuotaRepository) GetOrgQuotaQuantity(quotaID uuid.UUID, resourceID uuid.UUID) (uint, error) {
	var total uint
	err := r.db.
		Model(&models.ResourceQuantity{}).
		Joins("JOIN resource_properties rp ON rp.id = resource_quantities.resource_prop_id").
		Where("organization_quota_id = ? AND rp.resource_id = ?", quotaID, resourceID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *QuotaRepository) CreateProjectQuota(quota *models.ProjectQuota) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) GetProjectQuotaByProjectID(projectID uuid.UUID) ([]models.ProjectQuota, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.Debug().Preload("Resources.ResourceProp").Where("project_id = ?", projectID).
		Find(&projectQuotas).Error
	if err != nil {
		return nil, err
	}
	return projectQuotas, nil
}

func (r *QuotaRepository) GetProjectQuotaByID(id uuid.UUID) (*models.ProjectQuota, error) {
	var projectQuota models.ProjectQuota
	err := r.db.Preload("Resources.ResourceProp").First(&projectQuota, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &projectQuota, nil
}

func (r *QuotaRepository) CreateNamespaceQuota(quota *models.NamespaceQuota) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) GetNamespaceQuotaByNamespaceID(namespaceID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota

	err := r.db.Debug().
		Table("namespace_quota nq").
		Joins("JOIN namespace_quotas nqs ON nqs.namespace_quota_id = nq.id").
		Preload("Resources.ResourceProp").
		Preload("ResourcePool.Organization").
		Where("nqs.namespace_id = ?", namespaceID).
		Find(&namespaceQuotas).Error

	if err != nil {
		return nil, err
	}

	return namespaceQuotas, nil
}

func (r *QuotaRepository) CreateResourceProperty(resourceProperty *models.ResourceProperty) error {
	return r.db.Create(resourceProperty).Error
}

func (r *QuotaRepository) CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error {
	return r.db.Create(resourceQuantity).Error
}

func (r *QuotaRepository) GetNamespaceQuotaByID(id uuid.UUID) (*models.NamespaceQuota, error) {
	var namespaceQuota models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProp").First(&namespaceQuota, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &namespaceQuota, nil
}

func (r *QuotaRepository) GetOrganization(orgID uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.Preload("ResourcePools").First(&organization, "id = ?", orgID).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *QuotaRepository) GetNamespaceQuotasByProjectID(projectID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota

	err := r.db.Debug().
		Preload("Resources.ResourceProp").
		Preload("ResourcePool.Organization").
		Where("project_id = ?", projectID).
		Find(&namespaceQuotas).Error

	if err != nil {
		return nil, err
	}

	return namespaceQuotas, nil
}

func (r *QuotaRepository) IsOrgQuotaExist(fromOrgID uuid.UUID, toOrgID uuid.UUID, nodeID uuid.UUID) (bool, error) {
	var orgQuota models.OrganizationQuota
	err := r.db.Where("from_org_id = ? AND to_org_id = ? AND node_id = ?", fromOrgID, toOrgID, nodeID).
		First(&orgQuota).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (r *QuotaRepository) IsProjectQuotaExist(projectID, nodeID uuid.UUID) (bool, error) {
	var projectQuota models.ProjectQuota
	err := r.db.Where("project_id = ? AND node_id = ?", projectID, nodeID).
		First(&projectQuota).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *QuotaRepository) IsNamespaceQuotaExists(namespaceID uuid.UUID, nodeID uuid.UUID) (bool, error) {
	var namespaceQuota models.NamespaceQuota
	err := r.db.Where("namespace_id = ? AND node_id = ?", namespaceID, nodeID).
		First(&namespaceQuota).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

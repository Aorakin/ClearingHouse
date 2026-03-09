package repository

import (
	"sort"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/pkg/enum"
	"github.com/google/uuid"
)

func (r *QuotaRepository) CreateOrgQuota(quota *models.OrganizationQuota) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) GetOrganizationByRelationship(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuota, error) {
	var organizations []models.OrganizationQuota
	err := r.db.Preload("Resources.ResourceProp").Where("from_org_id = ? AND to_org_id = ?", fromOrgID, toOrgID).
		Order("name").
		Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *QuotaRepository) GetOrganizationQuotasByOrgID(orgID uuid.UUID) ([]models.OrganizationQuota, error) {
	var organizations []models.OrganizationQuota
	err := r.db.
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Preload("FromOrg").
		Preload("ToOrg").
		Where("from_org_id = ? OR to_org_id = ?", orgID, orgID).
		Order("name").
		Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	for i := range organizations {
		sort.Slice(organizations[i].Resources, func(a, b int) bool {
			return enum.ResourceTypeOrder(organizations[i].Resources[a].ResourceProp.Resource.ResourceType.Name) <
				enum.ResourceTypeOrder(organizations[i].Resources[b].ResourceProp.Resource.ResourceType.Name)
		})
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

func (r *QuotaRepository) DeleteOrganizationQuotasByOrgID(orgID uuid.UUID) error {
	return r.db.Where("to_org_id = ?", orgID).Delete(&models.OrganizationQuota{}).Error
}

func (r *QuotaRepository) HasActiveQuotaUsage(orgID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.Table("resource_quantities rq").
		Joins("JOIN project_quota pq ON pq.id = rq.project_quota_id").
		Joins("JOIN organization_quota oq ON oq.id = pq.organization_quota_id").
		Where("oq.to_org_id = ? AND rq.quantity > 0", orgID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *QuotaRepository) DeleteOrganizationQuota(quotaID uuid.UUID) error {
	// GORM will perform soft delete automatically (sets deleted_at)
	return r.db.Delete(&models.OrganizationQuota{}, "id = ?", quotaID).Error
}

func (r *QuotaRepository) HasProjectQuotasByOrgQuotaID(orgQuotaID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.Table("project_quota").
		Where("organization_quota_id = ? AND deleted_at IS NULL", orgQuotaID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *QuotaRepository) HasActiveUsageByOrgQuotaID(orgQuotaID uuid.UUID) (bool, error) {
	var count int64

	// Check if there are any resource quantities allocated to projects from this org quota
	// (not just the quantities defined in the org quota itself)
	err := r.db.Table("resource_quantities rq").
		Joins("JOIN project_quota pq ON pq.id = rq.project_quota_id").
		Where("pq.organization_quota_id = ? AND pq.deleted_at IS NULL AND rq.quantity > 0", orgQuotaID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

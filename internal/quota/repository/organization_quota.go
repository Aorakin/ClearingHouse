package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

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

func (r *QuotaRepository) GetOrganizationQuotasByOrgID(orgID uuid.UUID) ([]models.OrganizationQuota, error) {
	var organizations []models.OrganizationQuota
	err := r.db.
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Preload("FromOrg").
		Preload("ToOrg").
		Where("from_org_id = ? OR to_org_id = ?", orgID, orgID).
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

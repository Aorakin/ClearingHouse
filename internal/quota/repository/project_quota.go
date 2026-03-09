package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func (r *QuotaRepository) CreateProjectQuota(projectQuota *models.ProjectQuota) error {
	return r.db.Create(projectQuota).Error
}

func (r *QuotaRepository) GetProjectQuotaByProjectID(projectID uuid.UUID) ([]models.ProjectQuota, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.Debug().Preload("Resources.ResourceProp.Resource").Preload("Resources.ResourceProp.Resource.ResourceType").Where("project_id = ?", projectID).
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

func (r *QuotaRepository) DeleteProjectQuota(quotaID uuid.UUID) error {
	return r.db.Delete(&models.ProjectQuota{}, "id = ?", quotaID).Error
}

func (r *QuotaRepository) GetProjectQuotasByOrgQuotaID(orgQuotaID uuid.UUID) ([]models.ProjectQuota, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.Where("organization_quota_id = ?", orgQuotaID).Find(&projectQuotas).Error
	if err != nil {
		return nil, err
	}
	return projectQuotas, nil
}

func (r *QuotaRepository) DeleteResourceQuantitiesByProjectQuotaID(projectQuotaID uuid.UUID) error {
	return r.db.Where("project_quota_id = ?", projectQuotaID).Delete(&models.ResourceQuantity{}).Error
}

func (r *QuotaRepository) HasNamespaceQuotasByProjectQuotaID(projectQuotaID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.NamespaceQuota{}).
		Where("project_quota_id = ?", projectQuotaID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

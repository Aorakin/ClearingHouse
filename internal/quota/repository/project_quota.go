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

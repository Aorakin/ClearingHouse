package repository

import (
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func (r *QuotaRepository) CreateNamespaceQuotaTemplate(template *models.NamespaceQuotaTemplate) error {
	return r.db.Create(template).Error
}

func (r *QuotaRepository) GetNamespaceQuotaTemplateByID(templateID uuid.UUID) (*models.NamespaceQuotaTemplate, error) {
	var template models.NamespaceQuotaTemplate
	if err := r.db.Where("id = ?", templateID).First(&template).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *QuotaRepository) GetNamespaceQuotasByIDs(quotaIDs []uuid.UUID, projectID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota
	err := r.db.Where("id IN ? and project_id = ?", quotaIDs, projectID).
		Find(&namespaceQuotas).Error
	if err != nil {
		return nil, err
	}

	if len(namespaceQuotas) != len(quotaIDs) {
		return nil, fmt.Errorf("expected %d quotas but found %d", len(quotaIDs), len(namespaceQuotas))
	}

	return namespaceQuotas, nil
}

func (r *QuotaRepository) AssignQuotaToNamespace(namespaceID uuid.UUID, quotaTemplateID uuid.UUID) error {
	var namespace models.Namespace
	if err := r.db.First(&namespace, "id = ?", namespaceID).Error; err != nil {
		return err
	}
	namespace.QuotaTemplateID = &quotaTemplateID
	return r.db.Save(&namespace).Error
}

func (r *QuotaRepository) IsAssigned(namespaceID uuid.UUID, quotaID uuid.UUID) (bool, error) {
	var namespace models.Namespace
	if err := r.db.First(&namespace, "id = ?", namespaceID).Error; err != nil {
		return false, err
	}

	if namespace.QuotaTemplateID == nil {
		return false, nil
	}

	var quota models.NamespaceQuota
	if err := r.db.First(&quota, "id = ?", quotaID).Error; err != nil {
		return false, err
	}

	if quota.TemplateID == nil {
		return false, nil
	}

	return *quota.TemplateID == *namespace.QuotaTemplateID, nil
}

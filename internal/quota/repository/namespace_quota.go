package repository

import (
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func (r *QuotaRepository) CreateNamespaceQuotaTemplate(template *models.NamespaceQuotaTemplate) error {
	return r.db.Create(template).Error
}

func (r *QuotaRepository) AddQuotasToTemplate(templateID uuid.UUID, quotaIDs []uuid.UUID) error {
	var template models.NamespaceQuotaTemplate
	if err := r.db.First(&template, "id = ?", templateID).Error; err != nil {
		return err
	}

	var quotas []models.NamespaceQuota
	if err := r.db.Where("id IN ?", quotaIDs).Find(&quotas).Error; err != nil {
		return err
	}

	if len(quotas) != len(quotaIDs) {
		return fmt.Errorf("expected %d quotas but found %d", len(quotaIDs), len(quotas))
	}

	return r.db.Model(&template).Association("Quotas").Append(&quotas)
}

func (r *QuotaRepository) RemoveQuotasFromTemplate(templateID uuid.UUID, quotaIDs []uuid.UUID) error {
	var template models.NamespaceQuotaTemplate
	if err := r.db.First(&template, "id = ?", templateID).Error; err != nil {
		return err
	}

	var quotas []models.NamespaceQuota
	if err := r.db.Where("id IN ?", quotaIDs).Find(&quotas).Error; err != nil {
		return err
	}

	return r.db.Model(&template).Association("Quotas").Delete(&quotas)
}

func (r *QuotaRepository) GetNamespaceQuotaTemplateByID(templateID uuid.UUID) (*models.NamespaceQuotaTemplate, error) {
	var template models.NamespaceQuotaTemplate
	if err := r.db.Preload("Quotas").Where("id = ?", templateID).First(&template).Error; err != nil {
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

	// Check if the quota belongs to the template assigned to this namespace
	var count int64
	err := r.db.Table("namespace_quota_templates").
		Where("namespace_quota_template_id = ? AND namespace_quota_id = ?", *namespace.QuotaTemplateID, quotaID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

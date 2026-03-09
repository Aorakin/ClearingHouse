package repository

import (
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func (r *QuotaRepository) CreateNamespaceQuotaTemplate(template *models.NamespaceQuotaTemplate) error {
	return r.db.Create(template).Error
}

func (r *QuotaRepository) UpdateNamespaceQuotaTemplate(templateID uuid.UUID, name, description string) error {
	updates := make(map[string]interface{})
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&models.NamespaceQuotaTemplate{}).Where("id = ?", templateID).Updates(updates).Error
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
	if err := r.db.Preload("Quotas.Resources.ResourceProp.Resource.ResourceType").Where("id = ?", templateID).First(&template).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *QuotaRepository) GetNamespaceQuotaTemplatesByProjectID(projectID uuid.UUID) ([]models.NamespaceQuotaTemplate, error) {
	var templates []models.NamespaceQuotaTemplate
	if err := r.db.Preload("Quotas.Resources.ResourceProp.Resource.ResourceType").Where("project_id = ?", projectID).Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
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

func (r *QuotaRepository) UnassignQuotaTemplateFromNamespace(namespaceID uuid.UUID) error {
	return r.db.Model(&models.Namespace{}).Where("id = ?", namespaceID).Update("quota_template_id", nil).Error
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
	err := r.db.Table("namespace_quota_template_relations").
		Where("namespace_quota_template_id = ? AND namespace_quota_id = ?", *namespace.QuotaTemplateID, quotaID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *QuotaRepository) GetNamespaceQuotasByProjectID(projectID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota

	err := r.db.
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Preload("Node.ResourcePool.Organization").
		Where("project_id = ?", projectID).
		Find(&namespaceQuotas).Error

	if err != nil {
		return nil, err
	}

	return namespaceQuotas, nil
}

func (r *QuotaRepository) CreateNamespaceQuota(quota *models.NamespaceQuota) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) GetNamespaceQuotaByID(quotaID uuid.UUID) (*models.NamespaceQuota, error) {
	var namespaceQuota models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProp").First(&namespaceQuota, "id = ?", quotaID).Error
	if err != nil {
		return nil, err
	}
	return &namespaceQuota, nil
}

func (r *QuotaRepository) GetNamespaceQuotaByNamespaceID(namespaceID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota

	var namespace models.Namespace
	if err := r.db.First(&namespace, "id = ?", namespaceID).Error; err != nil {
		return nil, err
	}

	if namespace.QuotaTemplateID == nil {
		return namespaceQuotas, nil
	}

	err := r.db.
		Joins("JOIN namespace_quota_template_relations nqt ON nqt.namespace_quota_id = namespace_quota.id").
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Preload("Node.ResourcePool.Organization").
		Where("nqt.namespace_quota_template_id = ?", *namespace.QuotaTemplateID).
		Find(&namespaceQuotas).Error

	if err != nil {
		return nil, err
	}

	return namespaceQuotas, nil
}

func (r *QuotaRepository) DeleteNamespaceQuota(quotaID uuid.UUID) error {
	return r.db.Delete(&models.NamespaceQuota{}, "id = ?", quotaID).Error
}

func (r *QuotaRepository) GetNamespaceQuotasByProjectQuotaID(projectQuotaID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota
	err := r.db.Where("project_quota_id = ?", projectQuotaID).Find(&namespaceQuotas).Error
	if err != nil {
		return nil, err
	}
	return namespaceQuotas, nil
}

// RemoveNamespaceQuotaFromAllTemplates removes the quota from all templates via the junction table.
// This does NOT delete the templates themselves since it is a many2many relationship.
func (r *QuotaRepository) RemoveNamespaceQuotaFromAllTemplates(quotaID uuid.UUID) error {
	return r.db.Exec("DELETE FROM namespace_quota_template_relations WHERE namespace_quota_id = ?", quotaID).Error
}

// UnassignAllNamespacesFromTemplate sets quota_template_id = NULL for every namespace using the given template.
func (r *QuotaRepository) UnassignAllNamespacesFromTemplate(templateID uuid.UUID) error {
	return r.db.Model(&models.Namespace{}).Where("quota_template_id = ?", templateID).Update("quota_template_id", nil).Error
}

// DeleteTemplateQuotaRelations removes all quota associations from a template via the junction table.
// This does NOT delete the NamespaceQuotas themselves.
func (r *QuotaRepository) DeleteTemplateQuotaRelations(templateID uuid.UUID) error {
	return r.db.Exec("DELETE FROM namespace_quota_template_relations WHERE namespace_quota_template_id = ?", templateID).Error
}

func (r *QuotaRepository) UpdateNamespaceQuota(quotaID uuid.UUID, name, description string) error {
	updates := make(map[string]interface{})
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&models.NamespaceQuota{}).Where("id = ?", quotaID).Updates(updates).Error
}

func (r *QuotaRepository) UpdateResourceQuantity(quantityID uuid.UUID, quantity uint) error {
	return r.db.Model(&models.ResourceQuantity{}).Where("id = ?", quantityID).Update("quantity", quantity).Error
}

func (r *QuotaRepository) DeleteResourceQuantitiesByNamespaceQuotaID(namespaceQuotaID uuid.UUID) error {
	return r.db.Where("namespace_quota_id = ?", namespaceQuotaID).Delete(&models.ResourceQuantity{}).Error
}

func (r *QuotaRepository) GetResourceQuantitiesByNamespaceQuotaID(namespaceQuotaID uuid.UUID) ([]models.ResourceQuantity, error) {
	var quantities []models.ResourceQuantity
	err := r.db.Preload("ResourceProp").Where("namespace_quota_id = ?", namespaceQuotaID).Find(&quantities).Error
	if err != nil {
		return nil, err
	}
	return quantities, nil
}

func (r *QuotaRepository) HasQuotaTemplatesByNamespaceQuotaID(namespaceQuotaID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Table("namespace_quota_template_relations nqtr").
		Joins("INNER JOIN namespace_quota_templates nqt ON nqt.id = nqtr.namespace_quota_template_id").
		Where("nqtr.namespace_quota_id = ? AND nqt.deleted_at IS NULL", namespaceQuotaID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *QuotaRepository) DeleteNamespaceQuotaTemplate(templateID uuid.UUID) error {
	return r.db.Delete(&models.NamespaceQuotaTemplate{}, "id = ?", templateID).Error
}

func (r *QuotaRepository) HasNamespacesUsingTemplate(templateID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Namespace{}).
		Where("quota_template_id = ? AND deleted_at IS NULL", templateID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

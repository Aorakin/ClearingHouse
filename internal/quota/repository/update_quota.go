package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func (r *QuotaRepository) UpdateOrganizationQuota(quotaID uuid.UUID, name, description string) error {
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

	return r.db.Model(&models.OrganizationQuota{}).Where("id = ?", quotaID).Updates(updates).Error
}

func (r *QuotaRepository) UpdateProjectQuota(quotaID uuid.UUID, name, description string) error {
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

	return r.db.Model(&models.ProjectQuota{}).Where("id = ?", quotaID).Updates(updates).Error
}

func (r *QuotaRepository) UpdateResourceProperty(propID uuid.UUID, price float32, maxDuration uint) error {
	return r.db.Model(&models.ResourceProperty{}).Where("id = ?", propID).Updates(map[string]interface{}{
		"price":        price,
		"max_duration": maxDuration,
	}).Error
}

func (r *QuotaRepository) GetResourceQuantitiesByOrgQuotaID(orgQuotaID uuid.UUID) ([]models.ResourceQuantity, error) {
	var quantities []models.ResourceQuantity
	err := r.db.Preload("ResourceProp").Where("organization_quota_id = ?", orgQuotaID).Order("created_at").Find(&quantities).Error
	if err != nil {
		return nil, err
	}
	return quantities, nil
}

func (r *QuotaRepository) GetResourceQuantitiesByProjectQuotaID(projectQuotaID uuid.UUID) ([]models.ResourceQuantity, error) {
	var quantities []models.ResourceQuantity
	err := r.db.Preload("ResourceProp").Where("project_quota_id = ?", projectQuotaID).Order("created_at").Find(&quantities).Error
	if err != nil {
		return nil, err
	}
	return quantities, nil
}

func (r *QuotaRepository) GetProjectQuotasByOrgQuotaID(orgQuotaID uuid.UUID) ([]models.ProjectQuota, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.Preload("Resources.ResourceProp").Where("organization_quota_id = ?", orgQuotaID).Find(&projectQuotas).Error
	if err != nil {
		return nil, err
	}
	return projectQuotas, nil
}

func (r *QuotaRepository) GetNamespaceQuotasByProjectQuotaID(projectQuotaID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProp").Where("project_quota_id = ?", projectQuotaID).Find(&namespaceQuotas).Error
	if err != nil {
		return nil, err
	}
	return namespaceQuotas, nil
}

func (r *QuotaRepository) SoftDeleteResourceQuantitiesByOrgQuotaID(orgQuotaID uuid.UUID) error {
	return r.db.Where("organization_quota_id = ?", orgQuotaID).Delete(&models.ResourceQuantity{}).Error
}

func (r *QuotaRepository) SoftDeleteResourceQuantitiesByProjectQuotaID(projectQuotaID uuid.UUID) error {
	return r.db.Where("project_quota_id = ?", projectQuotaID).Delete(&models.ResourceQuantity{}).Error
}

func (r *QuotaRepository) SoftDeleteProjectQuotasByOrgQuotaID(orgQuotaID uuid.UUID) ([]uuid.UUID, error) {
	var projectQuotas []models.ProjectQuota
	if err := r.db.Where("organization_quota_id = ?", orgQuotaID).Find(&projectQuotas).Error; err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for _, pq := range projectQuotas {
		ids = append(ids, pq.ID)
	}

	if len(ids) > 0 {
		if err := r.db.Where("organization_quota_id = ?", orgQuotaID).Delete(&models.ProjectQuota{}).Error; err != nil {
			return nil, err
		}
	}

	return ids, nil
}

func (r *QuotaRepository) SoftDeleteNamespaceQuotasByProjectQuotaID(projectQuotaID uuid.UUID) ([]uuid.UUID, error) {
	var namespaceQuotas []models.NamespaceQuota
	if err := r.db.Where("project_quota_id = ?", projectQuotaID).Find(&namespaceQuotas).Error; err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for _, nq := range namespaceQuotas {
		ids = append(ids, nq.ID)
	}

	if len(ids) > 0 {
		if err := r.db.Where("project_quota_id = ?", projectQuotaID).Delete(&models.NamespaceQuota{}).Error; err != nil {
			return nil, err
		}
	}

	return ids, nil
}

func (r *QuotaRepository) UnassignQuotaTemplatesByNamespaceQuotaIDs(namespaceQuotaIDs []uuid.UUID) error {
	if len(namespaceQuotaIDs) == 0 {
		return nil
	}

	// Find all template IDs that reference these namespace quotas
	var templateIDs []uuid.UUID
	err := r.db.Table("namespace_quota_template_relations").
		Select("DISTINCT namespace_quota_template_id").
		Where("namespace_quota_id IN ?", namespaceQuotaIDs).
		Scan(&templateIDs).Error
	if err != nil {
		return err
	}

	// Remove the relations
	err = r.db.Exec(
		"DELETE FROM namespace_quota_template_relations WHERE namespace_quota_id IN ?",
		namespaceQuotaIDs,
	).Error
	if err != nil {
		return err
	}

	// For templates that now have no quotas, unassign them from namespaces
	for _, templateID := range templateIDs {
		var count int64
		if err := r.db.Table("namespace_quota_template_relations").
			Where("namespace_quota_template_id = ?", templateID).
			Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			// Template is now empty — unassign from all namespaces using it
			if err := r.db.Model(&models.Namespace{}).
				Where("quota_template_id = ?", templateID).
				Update("quota_template_id", nil).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

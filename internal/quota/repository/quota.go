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

func (r *QuotaRepository) CreateResourceProperty(resourceProperty *models.ResourceProperty) error {
	return r.db.Create(resourceProperty).Error
}

func (r *QuotaRepository) CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error {
	return r.db.Create(resourceQuantity).Error
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

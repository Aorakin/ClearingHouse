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
func (r *QuotaRepository) FindOrganizationQuotaGroup(fromOrgId uuid.UUID, toOrgId uuid.UUID) ([]models.OrganizationQuotaGroup, error) {
	var quotaGroups []models.OrganizationQuotaGroup
	err := r.db.Debug().Preload("Resources").Where("from_organization_id = ? AND to_organization_id = ?", fromOrgId, toOrgId).Find(&quotaGroups).Error
	if err != nil {
		return nil, err
	}
	return quotaGroups, nil
}

func (r *QuotaRepository) FindExistingOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID, poolID uuid.UUID) (*models.OrganizationQuotaGroup, error) {
	var quotaGroup models.OrganizationQuotaGroup
	err := r.db.Joins("JOIN resource_quantities rq ON rq.organization_quota_group_id = organization_quota_groups.id").
		Joins("JOIN resource_properties rp ON rq.resource_property_id = rp.id").
		Joins("JOIN resources r ON rp.resource_id = r.id").
		Where("from_organization_id = ? AND to_organization_id = ? AND r.resource_pool_id = ?", fromOrgID, toOrgID, poolID).
		First(&quotaGroup).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No existing quota group found
		}
		return nil, err // Other error
	}
	return &quotaGroup, nil
}

func (r *QuotaRepository) CreateOrganizationQuotaGroup(quota *models.OrganizationQuotaGroup) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error {
	return r.db.Create(resourceQuantity).Error
}

func (r *QuotaRepository) CreateResourceProperty(resourceProperty *models.ResourceProperty) error {
	return r.db.Create(resourceProperty).Error
}

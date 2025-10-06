package repository

import (
	"sort"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/ClearingHouse/pkg/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuotaRepository struct {
	db *gorm.DB
}

func NewQuotaRepository(db *gorm.DB) interfaces.QuotaRepository {
	return &QuotaRepository{db: db}
}

func (r *QuotaRepository) IsOrgQuotaExist(fromOrgID uuid.UUID, toOrgID uuid.UUID, poolID uuid.UUID) (bool, error) {
	var orgQuota models.OrganizationQuota
	err := r.db.Where("from_org_id = ? AND to_org_id = ? AND resource_pool_id = ?", fromOrgID, toOrgID, poolID).
		First(&orgQuota).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

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

func (r *QuotaRepository) CreateProjectQuota(quota *models.ProjectQuota) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) GetProjectQuotaByProjectID(projectID uuid.UUID) ([]models.ProjectQuota, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.Debug().Preload("Resources.ResourceProp").Where("project_id = ?", projectID).
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

func (r *QuotaRepository) CreateNamespaceQuota(quota *models.NamespaceQuota) error {
	return r.db.Create(quota).Error
}

func (r *QuotaRepository) GetNamespaceQuotaByNamespaceID(namespaceID uuid.UUID) ([]models.NamespaceQuota, error) {
	var namespaceQuotas []models.NamespaceQuota

	err := r.db.Debug().
		Table("namespace_quota nq").
		Joins("JOIN namespace_quotas nqs ON nqs.namespace_quota_id = nq.id").
		Preload("Resources.ResourceProp").
		Where("nqs.namespace_id = ?", namespaceID).
		Find(&namespaceQuotas).Error

	if err != nil {
		return nil, err
	}

	return namespaceQuotas, nil
}

func (r *QuotaRepository) CreateResourceProperty(resourceProperty *models.ResourceProperty) error {
	return r.db.Create(resourceProperty).Error
}

func (r *QuotaRepository) CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error {
	return r.db.Create(resourceQuantity).Error
}

func (r *QuotaRepository) GetNamespaceQuotaByID(id uuid.UUID) (*models.NamespaceQuota, error) {
	var namespaceQuota models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProp").First(&namespaceQuota, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &namespaceQuota, nil
}

func (r *QuotaRepository) AssignQuotaToNamespace(namespaceID uuid.UUID, namespaceQuotaID uuid.UUID) error {
	var nsQuota models.NamespaceQuota

	if err := r.db.Preload("Namespaces").First(&nsQuota, "id = ?", namespaceQuotaID).Error; err != nil {
		return err
	}

	for _, ns := range nsQuota.Namespaces {
		if ns.ID == namespaceID {
			return nil
		}
	}

	ns := models.Namespace{BaseModel: models.BaseModel{ID: namespaceID}}
	if err := r.db.Model(&nsQuota).Association("Namespaces").Append(&ns); err != nil {
		return err
	}

	return nil
}

func (r *QuotaRepository) IsAssigned(namespaceID uuid.UUID, quotaID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Table("namespace_quotas").
		Where("namespace_id = ? AND namespace_quota_id = ?", namespaceID, quotaID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *QuotaRepository) GetOrganization(orgID uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.Preload("ResourcePools").First(&organization, "id = ?", orgID).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *QuotaRepository) GetNamespaceUsageByType(namespaceID uuid.UUID, quotaID uuid.UUID) (*dtos.ResourceUsageResponse, error) {
	var tickets []models.Ticket

	if err := r.db.
		Where("tickets.namespace_id = ? and tickets.quota_id = ? and status IN ?", namespaceID, quotaID, enum.UsingStatuses).
		Preload("Resources.Resource.ResourceType").
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	typeAgg := make(map[string]dtos.ResourceUsage)
	for _, t := range tickets {
		for _, tr := range t.Resources {
			rt := tr.Resource.ResourceType
			rtID := rt.ID.String()

			if _, ok := typeAgg[rtID]; !ok {
				typeAgg[rtID] = dtos.ResourceUsage{
					TypeID: rtID,
					Type:   rt.Name,
					Usage:  0,
				}
			}
			tmp := typeAgg[rtID]
			tmp.Usage += float64(tr.Quantity)
			typeAgg[rtID] = tmp
		}
	}

	var result []dtos.ResourceUsage
	for _, v := range typeAgg {
		result = append(result, v)
	}

	return &dtos.ResourceUsageResponse{ResourceUsages: result}, nil
}

func (r *QuotaRepository) GetNamespaceQuotaByType(namespaceID uuid.UUID) (*dtos.ResourceQuotaResponse, error) {
	var quota models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProp.Resource.ResourceType").
		Joins("JOIN namespace_quotas nq ON nq.namespace_id = ?", namespaceID).
		First(&quota).Error
	if err != nil {
		return nil, err
	}

	typeAgg := make(map[string]dtos.ResourceQuota)
	for _, res := range quota.Resources {
		rt := res.ResourceProp.Resource.ResourceType
		rtID := rt.ID.String()

		if _, ok := typeAgg[rtID]; !ok {
			typeAgg[rtID] = dtos.ResourceQuota{
				TypeID: rtID,
				Type:   rt.Name,
				Quota:  0,
			}
		}
		tmp := typeAgg[rtID]
		tmp.Quota += float64(res.Quantity)
		typeAgg[rtID] = tmp
	}

	var result []dtos.ResourceQuota
	for _, v := range typeAgg {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Type < result[j].Type
	})

	return &dtos.ResourceQuotaResponse{ResourceQuotas: result}, nil
}

func (r *QuotaRepository) IsNamespaceQuotaExists(namespaceID uuid.UUID, resourcePoolID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Debug().Table("namespace_quotas").
		Joins("JOIN namespace_quota nq ON nq.id = namespace_quotas.namespace_quota_id").
		Where("namespace_id = ? AND resource_pool_id = ?", namespaceID, resourcePoolID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

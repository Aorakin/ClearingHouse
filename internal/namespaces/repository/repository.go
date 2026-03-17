package repository

import (
	"log"
	"sort"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/ClearingHouse/pkg/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NamespaceRepository struct {
	db *gorm.DB
}

func NewNamespaceRepository(db *gorm.DB) interfaces.NamespaceRepository {
	return &NamespaceRepository{
		db: db,
	}
}

func (r *NamespaceRepository) Create(namespace *models.Namespace) error {
	return r.db.Create(namespace).Error
}

func (r *NamespaceRepository) GetAll() ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Order("name").Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) GetNamespaceByID(namespaceID uuid.UUID) (*models.Namespace, error) {
	var namespace models.Namespace
	err := r.db.Preload("Members").Preload("Owner").First(&namespace, "id = ?", namespaceID).Error
	if err != nil {
		return nil, err
	}
	return &namespace, nil
}

func (r *NamespaceRepository) GetAllNamespacesByProjectID(projectID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.
		Joins("LEFT JOIN projects ON projects.id = namespaces.project_id").
		Where("namespaces.project_id = ? AND projects.deleted_at IS NULL", projectID).
		Preload("Owner").
		Preload("QuotaTemplate.Quotas.Resources.ResourceProp.Resource.ResourceType").
		Preload("Members").
		Order("namespaces.name DESC").
		Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) UpdateMembers(namespace *models.Namespace) error {
	return r.db.Model(namespace).Association("Members").Replace(namespace.Members)
}

func (r *NamespaceRepository) UpdateNamespace(namespace *models.Namespace) error {
	return r.db.Save(namespace).Error
}

func (r *NamespaceRepository) DeleteNamespace(namespaceID uuid.UUID) error {
	return r.db.Delete(&models.Namespace{}, "id = ?", namespaceID).Error
}

func (r *NamespaceRepository) HasActiveTicketsByNamespaceID(namespaceID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Ticket{}).
		Where("namespace_id = ? AND deleted_at IS NULL AND status NOT IN ?", namespaceID, []string{"failed", "cancelled", "expired", "stopped"}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *NamespaceRepository) GetAllNamespacesByUserID(userID uuid.UUID) ([]models.Namespace, error) {
	var user models.User
	if err := r.db.Debug().Preload("MemberNamespaces", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.MemberNamespaces, nil
}

func (r *NamespaceRepository) GetNamespaceQuotas(namespaceID uuid.UUID) ([]models.NamespaceQuota, error) {
	var quotas []models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProperties").Joins("JOIN namespace_quotas nq ON nq.namespace_id = ?", namespaceID).Order("namespace_quota.created_at").Find(&quotas).Error
	if err != nil {
		return nil, err
	}
	return quotas, nil
}

func (r *NamespaceRepository) GetNamespaceTickets(namespaceID, resourcePoolID, quotaID uuid.UUID) ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.db.Debug().
		Preload("Resources").
		Where("namespace_id = ? AND resource_pool_id = ? AND quota_id = ?", namespaceID, resourcePoolID, quotaID).
		Order("created_at DESC").
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *NamespaceRepository) GetAllNamespacesByProjectAndUserID(projectID uuid.UUID, userID uuid.UUID) ([]models.Namespace, error) {
	log.Println("")
	var namespaces []models.Namespace
	err := r.db.Joins("JOIN namespace_members nm ON nm.namespace_id = namespaces.id").
		Where("namespaces.project_id = ? AND nm.user_id = ?", projectID, userID).
		Order("namespaces.name").
		Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) GetNamespaceQuotaByType(namespaceID uuid.UUID) (*dtos.ResourceQuotaResponse, error) {
	var namespace models.Namespace
	if err := r.db.First(&namespace, "id = ?", namespaceID).Error; err != nil {
		return nil, err
	}

	if namespace.QuotaTemplateID == nil {
		return &dtos.ResourceQuotaResponse{ResourceQuotas: []dtos.ResourceQuota{}}, nil
	}

	var quotas []models.NamespaceQuota
	err := r.db.
		Table("namespace_quota").
		Joins("JOIN namespace_quota_template_relations nqt ON nqt.namespace_quota_id = namespace_quota.id").
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Where("nqt.namespace_quota_template_id = ?", *namespace.QuotaTemplateID).
		Find(&quotas).Error
	if err != nil {
		return nil, err
	}

	typeAgg := make(map[uuid.UUID]dtos.ResourceQuota)
	for _, quota := range quotas {
		for _, res := range quota.Resources {
			rt := res.ResourceProp.Resource.ResourceType
			rtID := rt.ID

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
	}

	var result []dtos.ResourceQuota
	for _, v := range typeAgg {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return enum.ResourceTypeOrder(result[i].Type) < enum.ResourceTypeOrder(result[j].Type)
	})

	return &dtos.ResourceQuotaResponse{ResourceQuotas: result}, nil
}

func (r *NamespaceRepository) GetNamespaceUsageByType(namespaceID uuid.UUID) (*dtos.ResourceUsageResponse, error) {
	var tickets []models.Ticket

	if err := r.db.
		Where("tickets.namespace_id = ? AND status IN ?", namespaceID, enum.UsingStatuses).
		Preload("Resources.Resource.ResourceType").
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	typeAgg := make(map[uuid.UUID]dtos.ResourceUsage)
	for _, t := range tickets {
		for _, tr := range t.Resources {
			rt := tr.Resource.ResourceType
			rtID := rt.ID

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

	sort.Slice(result, func(i, j int) bool {
		return enum.ResourceTypeOrder(result[i].Type) < enum.ResourceTypeOrder(result[j].Type)
	})

	return &dtos.ResourceUsageResponse{ResourceUsages: result}, nil
}

func (r *NamespaceRepository) GetPrivateNamespaceByUserID(userID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Where("owner_id = ?", userID).Order("name").Find(&namespaces).Error
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

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
	err := r.db.Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) GetNamespaceByID(namespaceID uuid.UUID) (*models.Namespace, error) {
	var namespace models.Namespace
	err := r.db.Preload("Members").First(&namespace, "id = ?", namespaceID).Error
	if err != nil {
		return nil, err
	}
	return &namespace, nil
}

func (r *NamespaceRepository) GetAllNamespacesByProjectID(projectID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Where("project_id = ?", projectID).Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) UpdateMembers(namespace *models.Namespace) error {
	return r.db.Model(namespace).Association("Members").Replace(namespace.Members)
}

func (r *NamespaceRepository) UpdateNamespace(namespace *models.Namespace) error {
	return r.db.Save(namespace).Error
}

func (r *NamespaceRepository) GetAllNamespacesByUserID(userID uuid.UUID) ([]models.Namespace, error) {
	var user models.User
	if err := r.db.Debug().Preload("MemberNamespaces").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.MemberNamespaces, nil
}

func (r *NamespaceRepository) GetNamespaceQuotas(namespaceID uuid.UUID) ([]models.NamespaceQuota, error) {
	var quotas []models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProperties").Joins("JOIN namespace_quotas nq ON nq.namespace_id = ?", namespaceID).Find(&quotas).Error
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
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *NamespaceRepository) GetAllNamespacesByProjectAndUserID(projectID, userID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Joins("JOIN namespace_members nm ON nm.namespace_id = namespaces.id").
		Where("namespaces.project_id = ? AND nm.user_id = ?", projectID, userID).
		Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) GetNamespaceQuotaByType(namespaceID uuid.UUID) (*dtos.ResourceQuotaResponse, error) {
	var quotas []models.NamespaceQuota
	err := r.db.Debug().
		Model(&models.NamespaceQuota{}).
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Joins("JOIN namespace_quotas nq ON nq.namespace_quota_id = namespace_quota.id").
		Where("nq.namespace_id = ?", namespaceID).
		Find(&quotas).Error
	if err != nil {
		return nil, err
	}

	log.Printf("Fetched %d quotas for namespace %s", len(quotas), namespaceID)
	log.Printf("Quotas: %+v", quotas)

	for _, q := range quotas {
		log.Printf("Quota ID: %s, Name: %s", q.ID, q.Name)
	}

	typeAgg := make(map[string]dtos.ResourceQuota)
	for _, quota := range quotas {
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

func (r *NamespaceRepository) GetNamespaceUsageByType(namespaceID uuid.UUID) (*dtos.ResourceUsageResponse, error) {
	var tickets []models.Ticket

	if err := r.db.
		Where("tickets.namespace_id = ? AND status IN ?", namespaceID, enum.UsingStatuses).
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

func (r *NamespaceRepository) GetPrivateNamespaceByUserID(userID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Where("owner_id = ?", userID).Find(&namespaces).Error
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

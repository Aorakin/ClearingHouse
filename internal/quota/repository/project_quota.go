package repository

import (
	"sort"

	"github.com/ClearingHouse/internal/models"
	namespaceDtos "github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/pkg/enum"
	"github.com/google/uuid"
)

func (r *QuotaRepository) CreateProjectQuota(projectQuota *models.ProjectQuota) error {
	return r.db.Create(projectQuota).Error
}

func (r *QuotaRepository) GetProjectQuotaByProjectID(projectID uuid.UUID) ([]models.ProjectQuota, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.Debug().Preload("Resources.ResourceProp.Resource").Preload("Resources.ResourceProp.Resource.ResourceType").Where("project_id = ?", projectID).
		Order("name").
		Find(&projectQuotas).Error
	if err != nil {
		return nil, err
	}
	for i := range projectQuotas {
		sort.Slice(projectQuotas[i].Resources, func(a, b int) bool {
			return enum.ResourceTypeOrder(projectQuotas[i].Resources[a].ResourceProp.Resource.ResourceType.Name) <
				enum.ResourceTypeOrder(projectQuotas[i].Resources[b].ResourceProp.Resource.ResourceType.Name)
		})
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

func (r *QuotaRepository) GetInternalProjectQuotasByOrgAndNode(orgID uuid.UUID, nodeID uuid.UUID) ([]models.ProjectQuota, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.
		Preload("Resources.ResourceProp").
		Where("organization_id = ? AND node_id = ? AND organization_quota_id IS NULL", orgID, nodeID).
		Find(&projectQuotas).Error
	if err != nil {
		return nil, err
	}
	return projectQuotas, nil
}

func (r *QuotaRepository) DeleteProjectQuota(quotaID uuid.UUID) error {
	return r.db.Delete(&models.ProjectQuota{}, "id = ?", quotaID).Error
}

func (r *QuotaRepository) GetProjectQuotaTotalByType(projectID uuid.UUID) (*namespaceDtos.ResourceQuotaResponse, error) {
	var projectQuotas []models.ProjectQuota
	err := r.db.
		Preload("Resources.ResourceProp.Resource").
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Where("project_id = ?", projectID).
		Find(&projectQuotas).Error
	if err != nil {
		return nil, err
	}

	typeAgg := make(map[uuid.UUID]namespaceDtos.ResourceQuota)
	for _, pq := range projectQuotas {
		for _, res := range pq.Resources {
			rt := res.ResourceProp.Resource.ResourceType
			rtID := rt.ID
			if _, ok := typeAgg[rtID]; !ok {
				typeAgg[rtID] = namespaceDtos.ResourceQuota{
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

	var result []namespaceDtos.ResourceQuota
	for _, v := range typeAgg {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return enum.ResourceTypeOrder(result[i].Type) < enum.ResourceTypeOrder(result[j].Type)
	})

	return &namespaceDtos.ResourceQuotaResponse{ResourceQuotas: result}, nil
}

func (r *QuotaRepository) HasNamespaceQuotasByProjectQuotaID(projectQuotaID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.NamespaceQuota{}).
		Where("project_quota_id = ?", projectQuotaID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

package repository

import (
	"log"
	"sort"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/pkg/enum"
	"github.com/google/uuid"
)

func (r *QuotaRepository) GetNamespaceUsageByType(namespaceID uuid.UUID, quotaID uuid.UUID) (*dtos.ResourceUsageResponse, error) {
	var tickets []models.Ticket

	if err := r.db.Debug().
		Where("tickets.namespace_id = ? and tickets.quota_id = ? and status IN ?", namespaceID, quotaID, enum.UsingStatuses).
		Preload("Resources.Resource.ResourceType").
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	log.Println("Tickets found:", len(tickets))

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
		return result[i].Type < result[j].Type
	})

	return &dtos.ResourceUsageResponse{ResourceUsages: result}, nil
}

func (r *QuotaRepository) GetNamespaceQuotaByType(namespaceID uuid.UUID) (*dtos.ResourceQuotaResponse, error) {
	var namespace models.Namespace
	err := r.db.First(&namespace, "id = ?", namespaceID).Error
	if err != nil {
		return nil, err
	}

	if namespace.QuotaTemplateID == nil {
		return &dtos.ResourceQuotaResponse{ResourceQuotas: []dtos.ResourceQuota{}}, nil
	}

	// Get all quotas belonging to this template
	var quotas []models.NamespaceQuota
	err = r.db.Preload("Resources.ResourceProp.Resource.ResourceType").
		Joins("JOIN namespace_quota_template_relations nqtr ON nqtr.namespace_quota_id = namespace_quota.id").
		Where("nqtr.namespace_quota_template_id = ?", *namespace.QuotaTemplateID).
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
		return result[i].Type < result[j].Type
	})

	return &dtos.ResourceQuotaResponse{ResourceQuotas: result}, nil
}

func (r *QuotaRepository) GetQuotaByType(quotaID uuid.UUID) (*dtos.ResourceQuotaResponse, error) {
	var quota models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProp.Resource.ResourceType").
		First(&quota, "namespace_quota.id = ?", quotaID).Error

	if err != nil {
		return nil, err
	}

	typeAgg := make(map[uuid.UUID]dtos.ResourceQuota)
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

	var result []dtos.ResourceQuota
	for _, v := range typeAgg {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Type < result[j].Type
	})

	return &dtos.ResourceQuotaResponse{ResourceQuotas: result}, nil
}

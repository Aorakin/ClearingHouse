package dtos

import "github.com/google/uuid"

type ResourceUsage struct {
	ResourcePoolID uuid.UUID             `json:"resource_pool_id"`
	Resources      []ResourceUsageDetail `json:"resources"`
}

type ResourceUsageDetail struct {
	ResourceID uuid.UUID `json:"resource_id"`
	Used       uint      `json:"used"`
	Total      uint      `json:"total"`
}

type QuotaUsageQuery struct {
	ResourcePoolID string `form:"resource_pool_id" binding:"required,uuid"`
	QuotaID        string `form:"quota_id" binding:"required,uuid"`
}

type QuotaUsageRequest struct {
	ResourcePoolID uuid.UUID `json:"resource_pool_id"`
	QuotaID        uuid.UUID `json:"quota_id"`
}

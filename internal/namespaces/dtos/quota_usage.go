package dtos

import "github.com/google/uuid"

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

type ResourceQuota struct {
	TypeID string  `json:"type_id"`
	Type   string  `json:"type"`
	Quota  float64 `json:"quota"`
}

type ResourceQuotaResponse struct {
	ResourceQuotas []ResourceQuota `json:"resource_quotas"`
}

type ResourceUsage struct {
	TypeID string  `json:"type_id"`
	Type   string  `json:"type"`
	Usage  float64 `json:"usage"`
}

type ResourceUsageResponse struct {
	ResourceUsages []ResourceUsage `json:"resource_usages"`
}

type NamespaceUsage struct {
	TypeID string  `json:"type_id"`
	Type   string  `json:"type"`
	Quota  float64 `json:"quota"`
	Usage  float64 `json:"usage"`
}

type NamespaceUsageResponse struct {
	Usage []NamespaceUsage `json:"usage"`
}

package dtos

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

type ProjectUsage struct {
	TypeID string  `json:"type_id"`
	Type   string  `json:"type"`
	Quota  float64 `json:"quota"`
	Usage  float64 `json:"usage"`
}

type ProjectUsageResponse struct {
	Usage []ProjectUsage `json:"usage"`
}

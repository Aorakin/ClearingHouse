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

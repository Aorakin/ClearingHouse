package dtos

import "github.com/google/uuid"

type UpdateOrganizationQuotaRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Resources   []OrganizationResources `json:"resources"`
}

type UpdateProjectQuotaRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Resources   []Resource `json:"resources"`
}

type UpdateInternalProjectQuotaRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Resources   []ResourceWithProperty `json:"resources"`
}

type UpdateResourcePropertyRequest struct {
	ResourcePropID uuid.UUID `json:"resource_property_id" binding:"required,uuid"`
	Price          float32   `json:"price" binding:"gte=0"`
	MaxDuration    uint      `json:"max_duration"`
}

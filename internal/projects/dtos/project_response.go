package dtos

import (
	"time"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/google/uuid"
)

type ProjectResponse struct {
	ID             uuid.UUID            `json:"id"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	OrganizationID uuid.UUID            `json:"organization_id"`
	Members        []models.User        `json:"members"`
	Admins         []models.User        `json:"admins"`
	ResourceQuotas []dtos.ResourceQuota `json:"resource_quotas"`
}

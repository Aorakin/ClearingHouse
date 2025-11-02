package dtos

import (
	"time"

	"github.com/ClearingHouse/internal/models"
)

type TicketResponse struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Status         string                  `json:"status"`
	StartTime      *time.Time              `json:"start_time"`
	EndTime        *time.Time              `json:"end_time"`
	CancelTime     *time.Time              `json:"cancel_time"`
	Duration       uint                    `json:"duration"`
	Price          float32                 `json:"price"`
	OwnerID        string                  `json:"owner_id"`
	NamespaceID    string                  `json:"namespace_id"`
	NamespaceName  string                  `json:"namespace_name"`
	ProjectID      string                  `json:"project_id"`
	ProjectName    string                  `json:"project_name"`
	ResourcePoolID string                  `json:"resource_pool_id"`
	QuotaID        string                  `json:"quota_id"`
	Resources      []models.TicketResource `json:"resources"`
	RedeemTimeout  uint                    `json:"redeem_timeout"`
}

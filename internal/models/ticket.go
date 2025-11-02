package models

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	BaseModel
	Name           string           `json:"name"`
	Status         string           `json:"status"`
	StartTime      *time.Time       `gorm:"type:timestamptz" json:"start_time"`
	EndTime        *time.Time       `gorm:"type:timestamptz" json:"end_time"`
	CancelTime     *time.Time       `gorm:"type:timestamptz" json:"cancel_time"`
	Duration       uint             `json:"duration"`
	Price          float32          `json:"price"`
	OwnerID        uuid.UUID        `gorm:"type:uuid;not null" json:"owner_id"`
	NamespaceID    uuid.UUID        `gorm:"type:uuid;not null" json:"namespace_id"`
	ResourcePoolID uuid.UUID        `gorm:"type:uuid;not null" json:"resource_pool_id"`
	QuotaID        uuid.UUID        `gorm:"type:uuid;not null" json:"quota_id"`
	Resources      []TicketResource `gorm:"foreignKey:TicketID" json:"resources"`
	Owner          User             `gorm:"foreignKey:OwnerID" json:"-"`
	Namespace      Namespace        `gorm:"foreignKey:NamespaceID" json:"-"`
	ResourcePool   ResourcePool     `gorm:"foreignKey:ResourcePoolID" json:"-"`
	RedeemTimeout  uint             `json:"redeem_timeout"` // in seconds
}

type TicketResource struct {
	ResourceID uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	TicketID   uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Quantity   uint      `json:"quantity"`
	Ticket     Ticket    `gorm:"foreignKey:TicketID" json:"-"`
	Resource   Resource  `gorm:"foreignKey:ResourceID" json:"-"`
}

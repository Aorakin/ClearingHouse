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
	Duration       uint             `json:"duration"`
	Price          float32          `json:"price"`
	OwnerID        uuid.UUID        `gorm:"type:uuid;not null" json:"owner_id"`
	NamespaceID    uuid.UUID        `gorm:"type:uuid;not null" json:"namespace_id"`
	ResourcePoolID uuid.UUID        `gorm:"type:uuid;not null" json:"resource_pool_id"`
	QuotaID        uuid.UUID        `gorm:"type:uuid;not null" json:"quota_id"`
	Resources      []TicketResource `gorm:"foreignKey:TicketID" json:"resources"`
	Owner          User             `gorm:"foreignKey:OwnerID" json:"-"`
	Namespace      Namespace        `gorm:"foreignKey:NamespaceID" json:"-"`
}

type TicketResource struct {
	ResourceID uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	TicketID   uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Quantity   uint      `json:"quantity"`
	Ticket     Ticket    `gorm:"foreignKey:TicketID" json:"-"`
	Resource   Resource  `gorm:"foreignKey:ResourceID" json:"-"`
}

type GliderTask struct {
	ID      uuid.UUID      `json:"id" validate:"required"`
	Tickets []GliderTicket `json:"tickets" validate:"required"`
	TaskID  string         `json:"task_id" validate:"required"`
}

type GliderTicket struct {
	ID                uuid.UUID  `json:"id" validate:"required"`
	NamespaceURN      string     `json:"namespace_urn" validate:"required"` // namespace.id
	GlideletURN       string     `json:"glidelet_urn" validate:"required"`  // resource_pool.URN
	Spec              GliderSpec `json:"spec" validate:"required"`
	ReferenceTicketID string     `json:"reference_ticket_id"`
	RedeemTimeout     uint       `json:"redeem_timeout" validate:"required"` // in seconds
	Lease             uint       `json:"lease" validate:"required"`          // in seconds
	CreatedAt         time.Time  `json:"created_at" validate:"required"`
}

type GliderSpec struct {
	Type      ResourceUnitType `json:"type" validate:"required"`
	PoolID    uuid.UUID        `json:"pool_id" validate:"required"`
	Resources []SpecResource   `json:"resource" validate:"required"`
}

type SpecResource struct {
	ResourceID uuid.UUID `json:"resource_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Quantity   string    `json:"quantity" validate:"required"`
	Unit       string    `json:"unit" validate:"required"`
}

type StatusTicket string

const (
	StatusReady    StatusTicket = "ready"
	StatusRedeemed StatusTicket = "redeemed"
	StatusCanceled StatusTicket = "canceled"
	StatusExpired  StatusTicket = "expired"
)

type ResourceUnitType string

const (
	ResourceUnitTypeCPU ResourceUnitType = "compute"
)

// type Task struct {
// 	ID          uuid.UUID     `json:"id"`
// 	NamespaceID uuid.UUID     `json:"namespace_id"`
// 	Namespace   Namespace     `gorm:"foreignKey:NamespaceID" json:"-"`
// 	Tickets     []TicketTable `gorm:"foreignKey:TaskID" json:"tickets"`
// }

// type TicketTable struct {
// 	ID         uuid.UUID  `json:"id"`
// 	TaskID     uuid.UUID  `json:"task_id"`
// 	RealTicket RealTicket `json:"real_ticket"`
// }

// // store text
// type RealTicket struct {
// }

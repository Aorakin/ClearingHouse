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
	Duration       float32          `json:"duration"`
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

type GliderTicket struct {
	ID                uuid.UUID    `json:"id" validate:"required"`
	NamespaceURN      string       `json:"namespace_urn" validate:"required"`
	GlideletURN       string       `json:"glidelet_urn" validate:"required"`
	Spec              []GliderSpec `json:"spec" validate:"required"`
	ReferenceTicketID string       `json:"reference_ticket_id"`
	RedeemTimeout     interface{}  `json:"redeem_timeout" validate:"required"`
	Lease             interface{}  `json:"lease" validate:"required"`
	Signature         interface{}  `json:"signature" validate:"required"`
}

type GliderSpec struct {
	ID        uuid.UUID        `json:"id" validate:"required"`
	Type      ResourceUnitType `json:"type" validate:"required"`
	PoolID    uuid.UUID        `json:"pool_id" validate:"required"`
	Resources []SpecResource   `json:"resource" validate:"required"`
}

type SpecResource struct {
	Name     string `json:"name" validate:"required"`
	Quantity string `json:"quantity" validate:"required"`
	Unit     string `json:"unit" validate:"required"`
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
	ResourceUnitTypeCPU    ResourceUnitType = "cpu"
	ResourceUnitTypeMemory ResourceUnitType = "memory"
	ResourceUnitTypeGPU    ResourceUnitType = "gpu"
	ResourceUnitTypeDisk   ResourceUnitType = "disk"
)

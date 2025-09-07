package models

import (
	"github.com/google/uuid"
)

type Ticket struct {
	BaseModel
	Name           string           `json:"name"`
	OwnerID        uuid.UUID        `gorm:"type:uuid;not null" json:"owner_id"`
	NamespaceID    uuid.UUID        `gorm:"type:uuid;not null" json:"namespace_id"`
	Status         string           `json:"status"`
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
}

package models

import (
	"github.com/google/uuid"
)

type Ticket struct {
	BaseModel
	Name        string           `json:"name"`
	OwnerID     uuid.UUID        `gorm:"type:uuid;not null" json:"owner_id"`
	NamespaceID uuid.UUID        `gorm:"type:uuid;not null" json:"namespace_id"`
	Owner       User             `gorm:"foreignKey:OwnerID" json:"owner"`
	Namespace   Namespace        `gorm:"foreignKey:NamespaceID" json:"namespace"`
	Resources   []TicketResource `gorm:"foreignKey:TicketID" json:"resources"`
	Status      string           `json:"status"`
}

type TicketResource struct {
	QuotaID    uuid.UUID      `gorm:"type:uuid;not null" json:"quota_id"`
	ResourceID uuid.UUID      `gorm:"type:uuid;not null" json:"resource_id"`
	TicketID   uuid.UUID      `gorm:"type:uuid;not null" json:"ticket_id"`
	Quantity   uint           `json:"quantity"`
	Quota      NamespaceQuota `gorm:"foreignKey:QuotaID" json:"-"`
	Ticket     Ticket         `gorm:"foreignKey:TicketID" json:"-"`
}

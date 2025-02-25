package models

import "github.com/google/uuid"

type Namespace struct {
	BaseModel
	Title     string    `json:"title"`
	ProjectID uuid.UUID `json:"project_id"`
	QuotaID   uuid.UUID `json:"quota_id"`
	Project   Project   `gorm:"foreignKey:ProjectID"  json:"-"`
	Quota     Quota     `gorm:"foreignKey:QuotaID"  json:"-"`
}

type Ticket struct {
	NamespaceID uuid.UUID `json:"namespace_id"`
	Namespace   Namespace `gorm:"foreignKey:NamespaceID"  json:"-"`
	RequesterID uuid.UUID `json:"user_id"`
	Requester   User      `gorm:"foreignKey:RequesterID"  json:"-"`
}

package models

import "github.com/google/uuid"

type Project struct {
	BaseModel
	Name           string         `json:"name" `
	Description    string         `json:"description"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE;" json:"-"`
	Namespaces     []Namespace    `gorm:"foreignKey:ProjectID" json:"-"`
	Quotas         []ProjectQuota `gorm:"foreignKey:ProjectID" json:"-"`
	Members        []User         `gorm:"many2many:project_members;constraint:OnDelete:CASCADE;" json:"members"`
	Admins         []User         `gorm:"many2many:project_admins;constraint:OnDelete:CASCADE;" json:"admins"`
}

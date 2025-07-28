package models

import "github.com/google/uuid"

type Project struct {
	BaseModel
	Name           string       `json:"name" `
	Description    string       `json:"description"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Namespaces     []Namespace  `gorm:"foreignKey:ProjectID" json:"namespaces"`
	ProjectMembers []User       `gorm:"many2many:project_members;" json:"project_members"`
	ProjectAdmins  []User       `gorm:"many2many:project_admins;" json:"project_admins"`
}

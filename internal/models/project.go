package models

import "github.com/google/uuid"

type Project struct {
	BaseModel
	Name           string              `json:"name" `
	Description    string              `json:"description"`
	OrganizationID uuid.UUID           `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization        `gorm:"foreignKey:OrganizationID" json:"-"`
	Namespaces     []Namespace         `gorm:"foreignKey:ProjectID" json:"namespaces"`
	Quotas         []ProjectQuotaGroup `gorm:"many2many:project_quotas;" json:"quotas"`
	Members        []User              `gorm:"many2many:project_members;" json:"project_members"`
	Admins         []User              `gorm:"many2many:project_admins;" json:"project_admins"`
}

// type Project struct {
// 	BaseModel
// 	Name           string         `json:"name" `
// 	Description    string         `json:"description"`
// 	OrganizationID uuid.UUID      `gorm:"type:uuid;not null" json:"organization_id"`
// 	Organization   Organization   `gorm:"foreignKey:OrganizationID" json:"-"`
// 	Namespaces     []Namespace    `gorm:"foreignKey:ProjectID" json:"namespaces"`
// 	Quotas         []ProjectQuota `gorm:"foreignKey:ProjectID" json:"quotas"`
// 	Members        []User         `gorm:"many2many:project_members;" json:"project_members"`
// 	Admins         []User         `gorm:"many2many:project_admins;" json:"project_admins"`
// }

package models

import "github.com/google/uuid"

type Project struct {
	BaseModel
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Quota       []Quota     `gorm:"foreignKey:ProjectID"  json:"-"`
	Namespaces  []Namespace `gorm:"foreignKey:ProjectID"  json:"-"`
	Owners      []*User     `gorm:"many2many:project_owners;" json:"owners"`
	Members     []*User     `gorm:"many2many:project_members;" json:"members"`
}

type Quota struct {
	BaseModel
	Title     string    `json:"title"`
	GPU       uint      `json:"GPU"`
	RAM       uint      `json:"RAM"`
	Storage   uint      `json:"storage"`
	ProjectID uuid.UUID `json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID"  json:"-"`
}

package models

type Organization struct {
	BaseModel
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Admins      []User    `gorm:"many2many:organization_admins;" json:"admins"`
	Members     []User    `gorm:"many2many:organization_members;" json:"members"`
	Projects    []Project `gorm:"foreignKey:OrganizationID" json:"projects"`
}

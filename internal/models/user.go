package models

type User struct {
	BaseModel
	Email               string         `json:"email"`
	FirstName           string         `json:"first_name"`
	LastName            string         `json:"last_name"`
	MemberOrganizations []Organization `gorm:"many2many:organization_members;" json:"-"`
	AdminOrganizations  []Organization `gorm:"many2many:organization_admins;" json:"-"`
	MemberProjects      []Project      `gorm:"many2many:project_members;" json:"-"`
	AdminProjects       []Project      `gorm:"many2many:project_admins;" json:"-"`
	MemberNamespaces    []Namespace    `gorm:"many2many:namespace_members;" json:"-"`
	Namespace           *Namespace     `gorm:"foreignKey:OwnerID" json:"namespace"`
}

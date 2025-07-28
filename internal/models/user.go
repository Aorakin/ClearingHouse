package models

type User struct {
	BaseModel
	Username            string         `json:"username"`
	Password            string         `json:"password"`
	Email               string         `json:"email"`
	MemberOrganizations []Organization `gorm:"many2many:organization_members;" json:"member_organizations"`
	AdminOrganizations  []Organization `gorm:"many2many:organization_admins;" json:"admin_organizations"`
	MemberProjects      []Project      `gorm:"many2many:project_members;" json:"member_projects"`
	AdminProjects       []Project      `gorm:"many2many:project_admins;" json:"admin_projects"`
	MemberNamespaces    []Namespace    `gorm:"many2many:namespace_members;" json:"member_namespaces"`
}

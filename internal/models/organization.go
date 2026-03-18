package models

type Organization struct {
	BaseModel
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Domain        string              `json:"domain" gorm:"index"`
	Admins        []User              `gorm:"many2many:organization_admins;constraint:OnDelete:CASCADE;" json:"admins"`
	Members       []User              `gorm:"many2many:organization_members;constraint:OnDelete:CASCADE;" json:"members"`
	Projects      []Project           `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE;" json:"projects"`
	ResourcePools []ResourcePool      `gorm:"foreignKey:OrganizationID" json:"-"`
	Quotas        []OrganizationQuota `gorm:"foreignKey:ToOrgID" json:"quotas"`
	GivenQuotas   []OrganizationQuota `gorm:"foreignKey:FromOrgID" json:"given_quotas"`
}

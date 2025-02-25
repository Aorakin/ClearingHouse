package models

type Project struct {
	BaseModel
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Namespaces  []Namespace `gorm:"foreignKey:ProjectID"  json:"-"`
	Owners      []*User     `gorm:"many2many:project_owners;" json:"owners"`
	Members     []*User     `gorm:"many2many:project_members;" json:"members"`
}

type Quota struct {
	BaseModel
	Title   string `json:"title"`
	GPU     uint   `json:"GPU"`
	RAM     uint   `json:"RAM"`
	Storage uint   `json:"storage"`
}

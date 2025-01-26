package dtos

type RegisterInput struct {
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}

type LoginInput struct {
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}

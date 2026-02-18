package dtos

type AssignSuperAdminRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RevokeSuperAdminRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SuperAdminResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

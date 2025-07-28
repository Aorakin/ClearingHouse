package dtos

type CreateOrganization struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type OrganizationURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateOrganization struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

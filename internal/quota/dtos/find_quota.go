package dtos

type FindOrganizationQuotaGroupRequest struct {
	FromOrganizationID string `form:"from" binding:"required,uuid"`
	ToOrganizationID   string `form:"to" binding:"required,uuid"`
}

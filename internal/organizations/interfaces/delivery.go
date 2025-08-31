package interfaces

import "github.com/gin-gonic/gin"

type OrganizationHandler interface {
	CreateOrganization() gin.HandlerFunc
	GetOrganizationByID() gin.HandlerFunc
	GetOrganizations() gin.HandlerFunc
	UpdateOrganization() gin.HandlerFunc
	DeleteOrganization() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
}

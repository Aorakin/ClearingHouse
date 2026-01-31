package interfaces

import "github.com/gin-gonic/gin"

type OrganizationHandler interface {
	CreateOrganization() gin.HandlerFunc
	GetOrganizationByID() gin.HandlerFunc
	GetAllOrganizations() gin.HandlerFunc
	UpdateOrganization() gin.HandlerFunc
	DeleteOrganization() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
	RemoveMembers() gin.HandlerFunc
	AddAdmins() gin.HandlerFunc
	RemoveAdmins() gin.HandlerFunc
	GetMembers() gin.HandlerFunc
	GetOrganizationMembers() gin.HandlerFunc
}

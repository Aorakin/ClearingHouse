package interfaces

import "github.com/gin-gonic/gin"

type OrganizationHandler interface {
	CreateOrganization() gin.HandlerFunc
	GetOrganizationByID() gin.HandlerFunc
	GetAllOrganizations() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
}

package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/gin-gonic/gin"
)

func MapOrganizationRoutes(orgGroup *gin.RouterGroup, orgHandler interfaces.OrganizationHandler) {
	orgGroup.GET("/", orgHandler.GetAllOrganizations())
	orgGroup.Use(middleware.AuthMiddleware())
	orgGroup.GET("/:id", orgHandler.GetOrganizationByID())
	orgGroup.POST("/", orgHandler.CreateOrganization())
	orgGroup.POST("/members", orgHandler.AddMembers())
}

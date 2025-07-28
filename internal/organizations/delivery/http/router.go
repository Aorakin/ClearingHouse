package http

import (
	"github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/gin-gonic/gin"
)

func MapOrganizationRoutes(orgGroup *gin.RouterGroup, orgHandler interfaces.OrganizationHandler) {
	orgGroup.POST("/", orgHandler.CreateOrganization())
	orgGroup.GET("/:id", orgHandler.GetOrganizationByID())
	orgGroup.GET("/", orgHandler.GetOrganizations())
	orgGroup.PUT("/:id", orgHandler.UpdateOrganization())
	orgGroup.DELETE("/:id", orgHandler.DeleteOrganization())
}

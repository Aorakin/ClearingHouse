package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/gin-gonic/gin"
)

func MapOrganizationRoutes(orgGroup *gin.RouterGroup, orgHandler interfaces.OrganizationHandler) {
	orgGroup.Use(middleware.AuthMiddleware())
	orgGroup.GET("/", orgHandler.GetAllOrganizations())
	orgGroup.GET("/:org-id", orgHandler.GetOrganizationByID())
	orgGroup.POST("/", middleware.SuperAdminMiddleware(), orgHandler.CreateOrganization())
	orgGroup.PUT("/:org-id", middleware.SuperAdminMiddleware(), orgHandler.UpdateOrganization())
	orgGroup.DELETE("/:org-id", middleware.SuperAdminMiddleware(), orgHandler.DeleteOrganization())
	orgGroup.POST("/members", orgHandler.AddMembers())
	orgGroup.POST("/rm-members", orgHandler.RemoveMembers())
	orgGroup.POST("/admins", orgHandler.AddAdmins())
	orgGroup.POST("/rm-admins", orgHandler.RemoveAdmins())
	orgGroup.GET("/members", orgHandler.GetMembers())
	orgGroup.GET("/:org-id/members", orgHandler.GetOrganizationMembers())
}

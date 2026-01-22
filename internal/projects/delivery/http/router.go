package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/gin-gonic/gin"
)

func MapProjectRoutes(projectGroup *gin.RouterGroup, projectHandler interfaces.ProjectHandler) {
	projectGroup.GET("/", projectHandler.GetAllProjects())
	projectGroup.Use(middleware.AuthMiddleware())
	projectGroup.GET("/all", projectHandler.GetAllUserProjects())
	projectGroup.GET("/organization/:orgId", projectHandler.GetProjectsByOrganizationID())
	projectGroup.GET("/:id", projectHandler.GetProject())
	projectGroup.GET("/:id/members", projectHandler.GetProjectMembers())
	projectGroup.GET("/:id/usage", projectHandler.GetProjectUsage())
	projectGroup.POST("/", projectHandler.CreateProject())
	projectGroup.POST("/members", projectHandler.AddMembers())
	projectGroup.POST("/rm-members", projectHandler.RemoveMembers())
	projectGroup.PUT("/", projectHandler.UpdateProject())
	projectGroup.DELETE("/:id", projectHandler.DeleteProject())
}

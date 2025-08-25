package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/gin-gonic/gin"
)

func MapProjectRoutes(projectGroup *gin.RouterGroup, projectHandler interfaces.ProjectHandler) {
	projectGroup.GET("/", projectHandler.GetAllProjects())
	projectGroup.Use(middleware.AuthMiddleware())
	projectGroup.POST("/", projectHandler.CreateProject())
}

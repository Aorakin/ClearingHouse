package http

import (
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/gin-gonic/gin"
)

func MapAnnouncementsRoutes(projectsGroup *gin.RouterGroup, projectsHandlers interfaces.ProjectsHandlers) {
	projectsGroup.GET("/", projectsHandlers.GetAll())
	projectsGroup.GET("/:id", projectsHandlers.GetByID())
	projectsGroup.POST("/", projectsHandlers.Create())
	projectsGroup.PUT("/:id", projectsHandlers.Update())
	projectsGroup.DELETE("/:id", projectsHandlers.Delete())
}

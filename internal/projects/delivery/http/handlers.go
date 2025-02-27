package http

import (
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectsUsecase interfaces.ProjectsUsecase
}

func NewProjectsHandler(projectsUsecase interfaces.ProjectsUsecase) interfaces.ProjectsHandlers {
	return &ProjectHandler{projectsUsecase: projectsUsecase}
}

func (h *ProjectHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *ProjectHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func (h *ProjectHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func (h *ProjectHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *ProjectHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

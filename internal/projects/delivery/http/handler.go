package http

import (
	"log"
	"net/http"

	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projUsecase interfaces.ProjectUsecase
}

func NewProjectHandler(projUsecase interfaces.ProjectUsecase) interfaces.ProjectHandler {
	return &ProjectHandler{
		projUsecase: projUsecase,
	}
}

func (h *ProjectHandler) CreateProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.CreateProjectRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.projUsecase.CreateProject(&request); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
	}
}

func (h *ProjectHandler) GetAllProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		projects, err := h.projUsecase.GetAllProjects()
		log.Println("Fetching all projects with total count: ", len(projects))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, projects)
	}
}

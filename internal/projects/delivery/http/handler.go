package http

import (
	"log"
	"net/http"

	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.CreateProjectRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

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

func (h *ProjectHandler) AddMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.AddMembersRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

		project, err := h.projUsecase.AddMembers(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

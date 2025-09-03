package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/ClearingHouse/internal/projects/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
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
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateProjectRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		err := h.projUsecase.CreateProject(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "project created successfully"})
	}
}

func (h *ProjectHandler) GetAllProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		projects, err := h.projUsecase.GetAllProjects()

		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, projects)
	}
}

func (h *ProjectHandler) AddMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.AddMembersRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err.Error())))
			return
		}

		project, err := h.projUsecase.AddMembers(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

func (h *ProjectHandler) GetAllUserProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projects, err := h.projUsecase.GetAllUserProjects(userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, projects)
	}
}

func (h *ProjectHandler) GetProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projectID := c.Param("id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		projectUUID := uuid.MustParse(projectID)
		if projectUUID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		project, err := h.projUsecase.GetProjectByID(projectUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

func (h *ProjectHandler) GetProjectQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projectID := c.Param("id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		projectUUID := uuid.MustParse(projectID)
		if projectUUID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		project, err := h.projUsecase.GetProjectByID(projectUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

func (h *ProjectHandler) GetProjectUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projectID := c.Param("id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		projectUUID := uuid.MustParse(projectID)
		if projectUUID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		project, err := h.projUsecase.GetProjectByID(projectUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

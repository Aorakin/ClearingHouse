package http

import (
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuotaHandler struct {
	quotaUsecase interfaces.QuotaUsecase
}

func NewQuotaHandler(quotaUsecase interfaces.QuotaUsecase) *QuotaHandler {
	return &QuotaHandler{
		quotaUsecase: quotaUsecase,
	}
}

func (h *QuotaHandler) CreateOrganizationQuotaGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(400, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.CreateOrganizationQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		request.Creator = userID

		quotaGroup, err := h.quotaUsecase.CreateOrganizationQuotaGroup(&request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, quotaGroup)
	}
}

func (h *QuotaHandler) FindOrganizationQuotaGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(400, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.FindOrganizationQuotaGroupRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fromUUID, err := uuid.Parse(request.FromOrganizationID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		toUUID, err := uuid.Parse(request.ToOrganizationID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		quotaGroups, err := h.quotaUsecase.FindOrganizationQuotaGroup(fromUUID, toUUID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, quotaGroups)
	}
}

func (h *QuotaHandler) CreateProjectQuotaGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(400, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.CreateProjectQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

		quotaGroup, err := h.quotaUsecase.CreateProjectQuotaGroup(&request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, quotaGroup)
	}
}

func (h *QuotaHandler) FindProjectQuotaGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		if projectID == "" {
			c.JSON(400, gin.H{"error": "Project ID is required"})
			return
		}

		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		quotaGroup, err := h.quotaUsecase.FindProjectQuotaGroup(projectUUID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, quotaGroup)
	}
}

func (h *QuotaHandler) AssignQuotaToNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.AssignQuotaToNamespaceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := h.quotaUsecase.AssignQuotaToNamespace(&request); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(204, nil)
	}
}

func (h *QuotaHandler) CreateNamespaceQuotaGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.CreateNamespaceQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		quotaGroup, err := h.quotaUsecase.CreateNamespaceQuotaGroup(&request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, quotaGroup)
	}
}

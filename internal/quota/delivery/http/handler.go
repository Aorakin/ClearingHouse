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
		var request dtos.CreateOrganizationQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		c.JSON(200, gin.H{"message": "Request received", "data": request})
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

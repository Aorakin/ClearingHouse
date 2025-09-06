package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuotaHandler struct {
	quotaUsecase interfaces.QuotaUsecase
}

func NewQuotaHandler(quotaUsecase interfaces.QuotaUsecase) *interfaces.QuotaHandler {
	return &QuotaHandler{
		quotaUsecase: quotaUsecase,
	}
}

func (h *QuotaHandler) CreateOrganizationQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateOrganizationQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateOrganizationQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *QuotaHandler) GetOrganizationQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.FindOrganizationQuotaGroupRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		fromUUID, err := uuid.Parse(request.FromOrganizationID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}
		toUUID, err := uuid.Parse(request.ToOrganizationID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quotas, err := h.quotaUsecase.GetOrganizationQuota(fromUUID, toUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quotas)
	}
}

func (h *QuotaHandler) CreateProjectQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateProjectQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

// func (h *QuotaHandler) CreateOrganizationQuotaGroup() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID := c.MustGet("userID").(uuid.UUID)
// 		if userID == uuid.Nil {
// 			c.JSON(400, gin.H{"error": "unauthorized"})
// 			return
// 		}

// 		var request dtos.CreateOrganizationQuotaRequest
// 		if err := c.ShouldBindJSON(&request); err != nil {
// 			c.JSON(400, gin.H{"error": "Invalid request"})
// 			return
// 		}

// 		request.Creator = userID

// 		quotaGroup, err := h.quotaUsecase.CreateOrganizationQuotaGroup(&request)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(201, quotaGroup)
// 	}
// }

// func (h *QuotaHandler) FindOrganizationQuotaGroup() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID := c.MustGet("userID").(uuid.UUID)
// 		if userID == uuid.Nil {
// 			c.JSON(400, gin.H{"error": "unauthorized"})
// 			return
// 		}

// 		var request dtos.FindOrganizationQuotaGroupRequest
// 		if err := c.ShouldBindQuery(&request); err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 		fromUUID, err := uuid.Parse(request.FromOrganizationID)
// 		if err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}
// 		toUUID, err := uuid.Parse(request.ToOrganizationID)
// 		if err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 		quotaGroups, err := h.quotaUsecase.FindOrganizationQuotaGroup(fromUUID, toUUID)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(200, quotaGroups)
// 	}
// }

// func (h *QuotaHandler) CreateProjectQuotaGroup() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID := c.MustGet("userID").(uuid.UUID)
// 		if userID == uuid.Nil {
// 			c.JSON(400, gin.H{"error": "unauthorized"})
// 			return
// 		}

// 		var request dtos.CreateProjectQuotaRequest
// 		if err := c.ShouldBindJSON(&request); err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 		request.Creator = userID

// 		quotaGroup, err := h.quotaUsecase.CreateProjectQuotaGroup(&request)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(201, quotaGroup)
// 	}
// }

// func (h *QuotaHandler) FindProjectQuotaGroup() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		projectID := c.Param("id")
// 		if projectID == "" {
// 			c.JSON(400, gin.H{"error": "Project ID is required"})
// 			return
// 		}

// 		projectUUID, err := uuid.Parse(projectID)
// 		if err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 		quotaGroup, err := h.quotaUsecase.FindProjectQuotaGroup(projectUUID)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(200, quotaGroup)
// 	}
// }

// func (h *QuotaHandler) FindNamespaceQuotaGroup() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		namespaceID := c.Param("id")
// 		if namespaceID == "" {
// 			c.JSON(400, gin.H{"error": "Namespace ID is required"})
// 			return
// 		}

// 		namespaceUUID, err := uuid.Parse(namespaceID)
// 		if err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 		quotaGroup, err := h.quotaUsecase.GetNamespaceQuotaGroup(namespaceUUID)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(200, quotaGroup)
// 	}
// }

// func (h *QuotaHandler) AssignQuotaToNamespace() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID := c.MustGet("userID").(uuid.UUID)
// 		if userID == uuid.Nil {
// 			c.JSON(400, gin.H{"error": "unauthorized"})
// 			return
// 		}

// 		var request dtos.AssignQuotaToNamespaceRequest
// 		if err := c.ShouldBindJSON(&request); err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 		request.Creator = userID

// 		if err := h.quotaUsecase.AssignQuotaToNamespace(&request); err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(204, nil)
// 	}
// }

// func (h *QuotaHandler) CreateNamespaceQuotaGroup() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var request dtos.CreateNamespaceQuotaRequest
// 		if err := c.ShouldBindJSON(&request); err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 		quotaGroup, err := h.quotaUsecase.CreateNamespaceQuotaGroup(&request)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(201, quotaGroup)
// 	}
// }

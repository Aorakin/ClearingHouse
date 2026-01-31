package http

import (
	"net/http"

	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *QuotaHandler) DeleteProjectQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		quotaID := c.Param("quota_id")
		if quotaID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Quota ID is required")))
			return
		}

		quotaUUID, err := uuid.Parse(quotaID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.quotaUsecase.DeleteProjectQuota(quotaUUID, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Project quota deleted successfully"})
	}
}

func (h *QuotaHandler) DeleteNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		quotaID := c.Param("quota_id")
		if quotaID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Quota ID is required")))
			return
		}

		quotaUUID, err := uuid.Parse(quotaID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.quotaUsecase.DeleteNamespaceQuota(quotaUUID, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Namespace quota deleted successfully"})
	}
}
func (h *QuotaHandler) DeleteNamespaceQuotaTemplate() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		templateID := c.Param("quota_template_id")
		if templateID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Template ID is required")))
			return
		}

		templateUUID, err := uuid.Parse(templateID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.quotaUsecase.DeleteNamespaceQuotaTemplate(templateUUID, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Namespace quota template deleted successfully"})
	}
}

func (h *QuotaHandler) UnassignQuotaTemplateFromNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		namespaceID := c.Param("namespace_id")
		if namespaceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Namespace ID is required")))
			return
		}

		namespaceUUID, err := uuid.Parse(namespaceID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.quotaUsecase.UnassignQuotaTemplateFromNamespace(namespaceUUID, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Quota template unassigned successfully"})
	}
}

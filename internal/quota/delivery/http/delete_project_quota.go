package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/quota/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectQuotaDeleteHandler struct {
	quotaUsecase interfaces.QuotaUsecase
}

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

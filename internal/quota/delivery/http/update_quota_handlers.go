package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/quota/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *QuotaHandler) UpdateOrganizationQuota() gin.HandlerFunc {
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

		var request dtos.UpdateOrganizationQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.UpdateOrganizationQuota(quotaUUID, &request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quota)
	}
}

func (h *QuotaHandler) UpdateProjectQuota() gin.HandlerFunc {
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

		var request dtos.UpdateProjectQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.UpdateProjectQuota(quotaUUID, &request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quota)
	}
}

func (h *QuotaHandler) UpdateInternalProjectQuota() gin.HandlerFunc {
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

		var request dtos.UpdateInternalProjectQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.UpdateInternalProjectQuota(quotaUUID, &request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quota)
	}
}

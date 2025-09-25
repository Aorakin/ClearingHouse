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

func NewQuotaHandler(quotaUsecase interfaces.QuotaUsecase) interfaces.QuotaHandler {
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

		var request dtos.CreateProjectQuotaRequest
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

func (h *QuotaHandler) GetProjectQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Project ID is required")))
			return
		}

		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.GetProjectQuota(projectUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quota)
	}
}

func (h *QuotaHandler) CreateNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateNamespaceQuotaRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateNamespaceQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *QuotaHandler) GetNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaceID := c.Param("id")
		if namespaceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Namespace ID is required")))
			return
		}

		namespaceUUID, err := uuid.Parse(namespaceID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.GetNamespaceQuota(namespaceUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quota)
	}
}

func (h *QuotaHandler) AssignQuotaToNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.AssignQuotaToNamespaceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.quotaUsecase.AssignQuotaToNamespace(&request, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *QuotaHandler) CreateOwnedProjectQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateOwnedProjectQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateOwnedProjectQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *QuotaHandler) GetUsage() gin.HandlerFunc {
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
		quotaUUID, err := uuid.Parse(quotaID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		usage, err := h.quotaUsecase.GetUsage(quotaUUID, namespaceUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, usage)
	}
}

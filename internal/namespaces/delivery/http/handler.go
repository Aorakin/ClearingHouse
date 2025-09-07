package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NamespaceHandler struct {
	namespaceUsecase interfaces.NamespaceUsecase
}

func NewNamespaceHandler(namespaceUsecase interfaces.NamespaceUsecase) interfaces.NamespaceHandler {
	return &NamespaceHandler{
		namespaceUsecase: namespaceUsecase,
	}
}

func (h *NamespaceHandler) CreateNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateNamespaceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		namespace, err := h.namespaceUsecase.CreateNamespace(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}
		c.JSON(http.StatusCreated, namespace)
	}
}

func (h *NamespaceHandler) GetAllNamespaces() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaces, err := h.namespaceUsecase.GetAllNamespaces()
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}
		c.JSON(http.StatusOK, namespaces)
	}
}

func (h *NamespaceHandler) AddMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.AddMembersRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		namespace, err := h.namespaceUsecase.AddMembers(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, namespace)
	}
}

func (h *NamespaceHandler) GetAllUserNamespaces() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projID := c.Param("id")
		if projID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		projectUUID := uuid.MustParse(projID)
		if projectUUID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid project ID")))
			return
		}

		namespaces, err := h.namespaceUsecase.GetAllUserNamespaces(projectUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, namespaces)
	}
}

func (h *NamespaceHandler) GetNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		namespaceID := c.Param("id")
		if namespaceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid namespace ID")))
			return
		}

		namespaceUUID := uuid.MustParse(namespaceID)
		if namespaceUUID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid namespace ID")))
			return
		}

		namespace, err := h.namespaceUsecase.GetNamespace(namespaceUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, namespace)
	}
}

func (h *NamespaceHandler) GetNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		namespaceID := c.Param("id")
		if namespaceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid namespace ID")))
			return
		}

		namespaceUUID := uuid.MustParse(namespaceID)
		if namespaceUUID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid namespace ID")))
			return
		}

		quotas, err := h.namespaceUsecase.GetNamespaceQuotas(namespaceUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quotas)
	}
}

func (h *NamespaceHandler) GetNamespaceUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		namespaceID := c.Param("id")
		if namespaceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid namespace ID")))
			return
		}

		namespaceUUID := uuid.MustParse(namespaceID)
		if namespaceUUID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid namespace ID")))
			return
		}
		var request dtos.QuotaUsageQuery
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		resourcePoolUUID, err := uuid.Parse(request.ResourcePoolID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid resource pool ID")))
			return
		}

		quotaUUID, err := uuid.Parse(request.QuotaID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid quota ID")))
			return
		}

		requestUUID := dtos.QuotaUsageRequest{
			ResourcePoolID: resourcePoolUUID,
			QuotaID:        quotaUUID,
		}

		usages, err := h.namespaceUsecase.GetNamespaceUsages(&requestUUID, namespaceUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, usages)
	}
}

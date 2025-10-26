package http

import (
	"net/http"

	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/ClearingHouse/internal/private_namespaces/dtos"
	"github.com/ClearingHouse/internal/private_namespaces/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PrivateNamespaceHandler struct {
	privateNamespaceUsecase interfaces.PrivateNamespaceUsecase
	namespaceUsecase        namespaceInterfaces.NamespaceUsecase
}

func NewPrivateNamespaceHandler(privateNamespaceUsecase interfaces.PrivateNamespaceUsecase, namespaceUsecase namespaceInterfaces.NamespaceUsecase) interfaces.PrivateNamespaceHandler {
	return &PrivateNamespaceHandler{privateNamespaceUsecase: privateNamespaceUsecase, namespaceUsecase: namespaceUsecase}
}

func (h *PrivateNamespaceHandler) GetPrivateNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		namespace, err := h.privateNamespaceUsecase.GetPrivateNamespaceByOwnerID(userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, namespace)
	}
}

func (h *PrivateNamespaceHandler) CreatePrivateNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreatePrivateNamespaceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid request body")))
			return
		}

		namespace, err := h.privateNamespaceUsecase.CreatePrivateNamespace(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, namespace)
	}
}

func (h *PrivateNamespaceHandler) CreateNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateNamespaceQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid request body")))
			return
		}

		quota, err := h.privateNamespaceUsecase.CreateNamespaceQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *PrivateNamespaceHandler) GetUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		namespace, err := h.privateNamespaceUsecase.GetPrivateNamespaceByOwnerID(userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		usage, err := h.namespaceUsecase.GetNamespaceUsages(namespace.ID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, usage)
	}
}

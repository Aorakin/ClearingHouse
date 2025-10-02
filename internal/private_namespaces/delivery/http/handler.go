package http

import (
	"github.com/ClearingHouse/internal/private_namespaces/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PrivateNamespaceHandler struct {
	privateNamespaceUsecase interfaces.PrivateNamespaceUsecase
}

func NewPrivateNamespaceHandler(privateNamespaceUsecase interfaces.PrivateNamespaceUsecase) interfaces.PrivateNamespaceHandler {
	return &PrivateNamespaceHandler{privateNamespaceUsecase: privateNamespaceUsecase}
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

		c.JSON(200, namespace)
	}
}

func (h *PrivateNamespaceHandler) CreatePrivateNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Create Private Namespace",
		})
	}
}

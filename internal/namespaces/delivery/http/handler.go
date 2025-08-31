package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NamespaceHandler struct {
	usecase interfaces.NamespaceUsecase
}

func NewNamespaceHandler(usecase interfaces.NamespaceUsecase) interfaces.NamespaceHandler {
	return &NamespaceHandler{
		usecase: usecase,
	}
}

func (h *NamespaceHandler) CreateNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.CreateNamespaceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

		namespace, err := h.usecase.CreateNamespace(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, namespace)
	}
}

func (h *NamespaceHandler) GetAllNamespaces() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaces, err := h.usecase.GetAllNamespaces()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, namespaces)
	}
}

func (h *NamespaceHandler) AddMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.AddMembersRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

		namespace, err := h.usecase.AddMembers(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, namespace)
	}
}

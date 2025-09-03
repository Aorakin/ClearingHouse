package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var request dtos.CreateNamespaceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

		namespace, err := h.namespaceUsecase.CreateNamespace(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, namespace)
	}
}

func (h *NamespaceHandler) GetAllNamespaces() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaces, err := h.namespaceUsecase.GetAllNamespaces()
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

		namespace, err := h.namespaceUsecase.AddMembers(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, namespace)
	}
}

func (h *NamespaceHandler) GetAllUserNamespaces() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		namespaces, err := h.namespaceUsecase.GetAllUserNamespaces(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Message()})
			return
		}

		c.JSON(http.StatusOK, namespaces)
	}
}

func (h *NamespaceHandler) GetNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		namespaceID := c.Param("id")
		if namespaceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid namespace ID"})
			return
		}

		namespaceUUID := uuid.MustParse(namespaceID)
		if namespaceUUID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid namespace ID"})
			return
		}

		namespace, err := h.namespaceUsecase.GetNamespace(namespaceUUID, userID)
		if err != nil {
			c.JSON(err.Status(), gin.H{"error": err.Message()})
			return
		}

		c.JSON(http.StatusOK, namespace)
	}
}

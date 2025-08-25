package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
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
		var namespace dtos.CreateNamespaceRequest
		if err := c.ShouldBindJSON(&namespace); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.usecase.CreateNamespace(&namespace); err != nil {
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

package http

import (
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
)

type NamespacesHandler struct {
	namespacesUsecase interfaces.NamespacesUsecase
}

func NewNamespacesHandler(namespacesUsecase interfaces.NamespacesUsecase) interfaces.NamespacesHandlers {
	return &NamespacesHandler{namespacesUsecase: namespacesUsecase}
}

func (h *NamespacesHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *NamespacesHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *NamespacesHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *NamespacesHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *NamespacesHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

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

func (h *NamespacesHandler) AddMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *NamespacesHandler) ChangeQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *NamespacesHandler) RequestTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *NamespacesHandler) TerminateTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

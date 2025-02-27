package http

import (
	"github.com/ClearingHouse/internal/quotas/interfaces"
	"github.com/gin-gonic/gin"
)

type QuotasHandler struct {
	quotasUsecase interfaces.QuotasUsecase
}

func NewQuotasHandler(quotasUsecase interfaces.QuotasUsecase) interfaces.QuotasHandlers {
	return &QuotasHandler{quotasUsecase: quotasUsecase}
}

func (h *QuotasHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *QuotasHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func (h *QuotasHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func (h *QuotasHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (h *QuotasHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

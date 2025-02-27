package http

import (
	"github.com/ClearingHouse/internal/quotas/interfaces"
	"github.com/gin-gonic/gin"
)

func MapAnnouncementsRoutes(quotasGroup *gin.RouterGroup, quotasHandlers interfaces.QuotasHandlers) {
	quotasGroup.GET("/", quotasHandlers.GetAll())
	quotasGroup.GET("/:id", quotasHandlers.GetByID())
	quotasGroup.POST("/", quotasHandlers.Create())
	quotasGroup.PUT("/:id", quotasHandlers.Update())
	quotasGroup.DELETE("/:id", quotasHandlers.Delete())
}

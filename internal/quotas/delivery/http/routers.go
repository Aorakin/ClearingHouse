package http

import (
	"github.com/ClearingHouse/internal/quotas/interfaces"
	"github.com/gin-gonic/gin"
)

func MapAnnouncementsRoutes(quotasGroup *gin.RouterGroup, quotasHandlers interfaces.QuotasHandlers) {
	quotasGroup.GET("/", quotasHandlers.GetAll())
	quotasGroup.POST("/", quotasHandlers.Create())
}

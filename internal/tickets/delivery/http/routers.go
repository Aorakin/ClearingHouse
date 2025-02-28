package http

import (
	"github.com/ClearingHouse/internal/tickets/interfaces"
	"github.com/gin-gonic/gin"
)

func MapAnnouncementsRoutes(ticketsGroup *gin.RouterGroup, ticketsHandlers interfaces.TicketsHandlers) {
	ticketsGroup.GET("/", ticketsHandlers.GetAll())
	ticketsGroup.GET("/:id", ticketsHandlers.GetByID())
	ticketsGroup.POST("/", ticketsHandlers.Create())
	ticketsGroup.PUT("/:id", ticketsHandlers.Update())
	ticketsGroup.DELETE("/:id", ticketsHandlers.Delete())
}

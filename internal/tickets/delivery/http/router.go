package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	"github.com/gin-gonic/gin"
)

func MapTicketRoutes(ticketGroup *gin.RouterGroup, ticketHandler interfaces.TicketHandler) {
	ticketGroup.PATCH("/start", ticketHandler.StartTicket())
	ticketGroup.PATCH("/stop", ticketHandler.StopTicket())
	ticketGroup.Use(middleware.AuthMiddleware())
	ticketGroup.POST("/", ticketHandler.CreateTicket())
	ticketGroup.GET("/:ticket_id", ticketHandler.GetTicket())
	ticketGroup.PATCH("/:ticket_id/cancel", ticketHandler.CancelTicket())
	ticketGroup.GET("/namespace/:namespace_id", ticketHandler.GetNamespaceTickets())
	ticketGroup.GET("/user", ticketHandler.GetUserTickets())
}

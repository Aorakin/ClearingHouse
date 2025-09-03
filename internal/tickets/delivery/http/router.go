package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	"github.com/gin-gonic/gin"
)

func MapTicketRoutes(ticketGroup *gin.RouterGroup, ticketHandler interfaces.TicketHandler) {
	ticketGroup.Use(middleware.AuthMiddleware())
	ticketGroup.POST("/", ticketHandler.CreateTicket())
}

// get all ticket in ns
// get ticket detail
// refund ticket
// create ticket
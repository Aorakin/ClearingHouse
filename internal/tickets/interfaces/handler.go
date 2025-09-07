package interfaces

import "github.com/gin-gonic/gin"

type TicketHandler interface {
	CreateTicket() gin.HandlerFunc
	GetNamespaceTickets() gin.HandlerFunc
}

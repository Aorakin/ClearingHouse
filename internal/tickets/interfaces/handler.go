package interfaces

import "github.com/gin-gonic/gin"

type TicketHandler interface {
	CreateTicket() gin.HandlerFunc
	GetNamespaceTickets() gin.HandlerFunc
	GetUserTickets() gin.HandlerFunc
	StartTicket() gin.HandlerFunc
	StopTicket() gin.HandlerFunc
}

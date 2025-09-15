package interfaces

import "github.com/gin-gonic/gin"

type TicketHandler interface {
	CreateTicket() gin.HandlerFunc
	GetNamespaceTickets() gin.HandlerFunc
	GetUserTickets() gin.HandlerFunc
	GetTicket() gin.HandlerFunc
	StartTicket() gin.HandlerFunc
	StopTicket() gin.HandlerFunc
	CancelTicket() gin.HandlerFunc
}

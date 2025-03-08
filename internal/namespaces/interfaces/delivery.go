package interfaces

import "github.com/gin-gonic/gin"

type NamespacesHandlers interface {
	GetAll() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	Create() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
	ChangeQuota() gin.HandlerFunc
	RequestTicket() gin.HandlerFunc
	TerminateTicket() gin.HandlerFunc
}

package interfaces

import "github.com/gin-gonic/gin"

type QuotasHandlers interface {
	GetAll() gin.HandlerFunc
	Create() gin.HandlerFunc
}

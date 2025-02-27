package interfaces

import "github.com/gin-gonic/gin"

type QuotasHandlers interface {
	GetAll() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

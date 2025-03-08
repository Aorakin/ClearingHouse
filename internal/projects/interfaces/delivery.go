package interfaces

import "github.com/gin-gonic/gin"

type ProjectsHandlers interface {
	GetAll() gin.HandlerFunc
	Create() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
	// GetByID() gin.HandlerFunc
	// Update() gin.HandlerFunc
	// Delete() gin.HandlerFunc
}

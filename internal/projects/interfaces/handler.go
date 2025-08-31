package interfaces

import "github.com/gin-gonic/gin"

type ProjectHandler interface {
	CreateProject() gin.HandlerFunc
	GetAllProjects() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
	// GetProject() gin.HandlerFunc
	// UpdateProject() gin.HandlerFunc
	// DeleteProject() gin.HandlerFunc
}

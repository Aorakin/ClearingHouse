package interfaces

import "github.com/gin-gonic/gin"

type ProjectHandler interface {
	CreateProject() gin.HandlerFunc
	GetAllProjects() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
	// GetProject() gin.HandlerFunc
	// UpdateProject() gin.HandlerFunc
	// DeleteProject() gin.HandlerFunc

	GetAllUserProjects() gin.HandlerFunc
	GetProject() gin.HandlerFunc
	GetProjectQuota() gin.HandlerFunc
	GetProjectUsage() gin.HandlerFunc

	// get all user's available projects
	// get project detail
	// get quota of project
	// get current usage of project

}

package interfaces

import "github.com/gin-gonic/gin"

type ResourceHandler interface {
	GetResource() gin.HandlerFunc
	GetResourceTypes() gin.HandlerFunc
	CreateResourcePool() gin.HandlerFunc
	CreateResourceType() gin.HandlerFunc
	CreateResource() gin.HandlerFunc
	UpdateResource() gin.HandlerFunc

	GetResourceProperty() gin.HandlerFunc
	GetResourcePool() gin.HandlerFunc
}

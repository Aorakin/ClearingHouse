package http

import (
	"github.com/ClearingHouse/internal/resources/interfaces"
	"github.com/gin-gonic/gin"
)

func MapResourceRoutes(resourcesGroup *gin.RouterGroup, resourceHandler interfaces.ResourceHandler) {
	resourcesGroup.GET("/org/:id", resourceHandler.GetResource())
	resourcesGroup.PATCH("/:id", resourceHandler.UpdateResource())
	resourcesGroup.GET("/type", resourceHandler.GetResourceTypes())
	resourcesGroup.POST("/pool", resourceHandler.CreateResourcePool())
	resourcesGroup.POST("/type", resourceHandler.CreateResourceType())
	resourcesGroup.POST("/", resourceHandler.CreateResource())

	resourcesGroup.GET("/:id", resourceHandler.GetResourceProperty())
	resourcesGroup.GET("/pool/:id", resourceHandler.GetResourcePool())
}

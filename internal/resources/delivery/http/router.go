package http

import (
	"github.com/ClearingHouse/internal/resources/interfaces"
	"github.com/gin-gonic/gin"
)

func MapResourceRoutes(resourcesGroup *gin.RouterGroup, resourceHandler interfaces.ResourceHandler) {
	resourcesGroup.GET("/org/:id", resourceHandler.GetResource())
	resourcesGroup.GET("/type", resourceHandler.GetResourceTypes())
	resourcesGroup.GET("/:id", resourceHandler.GetResourceProperty())
	resourcesGroup.GET("/node/:node_id", resourceHandler.GetResourceNode())
	resourcesGroup.GET("/pool/:id", resourceHandler.GetResourcePool())
	resourcesGroup.POST("/pool", resourceHandler.CreateResourcePool())
	resourcesGroup.POST("/type", resourceHandler.CreateResourceType())
	resourcesGroup.POST("/node", resourceHandler.CreateResourceNode())
	resourcesGroup.POST("/", resourceHandler.CreateResource())
	resourcesGroup.PATCH("/:id", resourceHandler.UpdateResource())

}

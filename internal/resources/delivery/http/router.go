package http

import (
	"github.com/ClearingHouse/internal/resources/interfaces"
	"github.com/gin-gonic/gin"
)

func MapResourceRoutes(resourcesGroup *gin.RouterGroup, resourceHandler interfaces.ResourceHandler) {
	resourcesGroup.GET("/org/:org_id", resourceHandler.GetResource())
	resourcesGroup.POST("/type", resourceHandler.CreateResourceType())
	resourcesGroup.GET("/type", resourceHandler.GetResourceTypes())

	resourcesGroup.POST("/pool", resourceHandler.CreateResourcePool())
	resourcesGroup.GET("/pool/:pool_id", resourceHandler.GetResourcePool())
	resourcesGroup.PATCH("/pool/:pool_id", resourceHandler.UpdateResourcePool())
	resourcesGroup.DELETE("/pool/:pool_id", resourceHandler.DeleteResourcePool())
	resourcesGroup.POST("/node", resourceHandler.CreateResourceNode())
	resourcesGroup.GET("/node/:node_id", resourceHandler.GetResourceNode())
	resourcesGroup.PATCH("/node/:node_id", resourceHandler.UpdateResourceNode())
	resourcesGroup.DELETE("/node/:node_id", resourceHandler.DeleteResourceNode())
	resourcesGroup.GET("/:resource_id", resourceHandler.GetResourceProperty())
	resourcesGroup.POST("/", resourceHandler.CreateResource())
	resourcesGroup.PATCH("/:resource_id", resourceHandler.UpdateResource())
	resourcesGroup.DELETE("/:resource_id", resourceHandler.DeleteResource())

}

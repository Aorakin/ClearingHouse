package http

import (
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
)

func MapAnnouncementsRoutes(namespacesGroup *gin.RouterGroup, namespacesHandlers interfaces.NamespacesHandlers) {
	namespacesGroup.GET("/", namespacesHandlers.GetAll())
	namespacesGroup.GET("/:id", namespacesHandlers.GetByID())
	namespacesGroup.POST("/", namespacesHandlers.Create())
	namespacesGroup.PUT("/:id", namespacesHandlers.Update())
	namespacesGroup.DELETE("/:id", namespacesHandlers.Delete())
}

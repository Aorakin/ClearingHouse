package http

import (
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
)

func MapAnnouncementsRoutes(namespacesGroup *gin.RouterGroup, namespacesHandlers interfaces.NamespacesHandlers) {
	namespacesGroup.GET("/all/:project_id", namespacesHandlers.GetAll())
	namespacesGroup.POST("/", namespacesHandlers.Create())
	namespacesGroup.GET("/:id", namespacesHandlers.GetByID())
	namespacesGroup.PATCH("/:id/quota", namespacesHandlers.ChangeQuota())
	namespacesGroup.POST("/:id/member", namespacesHandlers.AddMembers())
	namespacesGroup.POST("/:id/ticket", namespacesHandlers.RequestTicket())
	namespacesGroup.DELETE("/:id/ticket", namespacesHandlers.TerminateTicket())

	// namespacesGroup.PUT("/:id", namespacesHandlers.Update())
	// namespacesGroup.DELETE("/:id", namespacesHandlers.Delete())
}

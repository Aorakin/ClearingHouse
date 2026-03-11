package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
)

func MapNamespaceRoutes(namespaceGroup *gin.RouterGroup, namespaceHandler interfaces.NamespaceHandler) {
	namespaceGroup.Use(middleware.AuthMiddleware())
	namespaceGroup.GET("/", namespaceHandler.GetAllNamespaces())
	namespaceGroup.POST("/", namespaceHandler.CreateNamespace())
	namespaceGroup.GET("/all/:id", namespaceHandler.GetAllUserNamespaces())
	namespaceGroup.GET("/project/:projectId", namespaceHandler.GetNamespacesByProjectID())
	namespaceGroup.GET("/:id", namespaceHandler.GetNamespace())
	namespaceGroup.GET("/:id/usage", namespaceHandler.GetNamespaceUsage())
	namespaceGroup.POST("/members", namespaceHandler.AddMembers())
	namespaceGroup.POST("/rm-members", namespaceHandler.RemoveMembers())
	namespaceGroup.PUT("/", namespaceHandler.UpdateNamespace())
	namespaceGroup.DELETE("/:id", namespaceHandler.DeleteNamespace())
}

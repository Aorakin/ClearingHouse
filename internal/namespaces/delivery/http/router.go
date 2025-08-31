package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
)

func MapNamespaceRoutes(namespaceGroup *gin.RouterGroup, namespaceHandler interfaces.NamespaceHandler) {
	namespaceGroup.GET("/", namespaceHandler.GetAllNamespaces())
	namespaceGroup.Use(middleware.AuthMiddleware())
	namespaceGroup.POST("/", namespaceHandler.CreateNamespace())
	namespaceGroup.POST("/members", namespaceHandler.AddMembers())
}

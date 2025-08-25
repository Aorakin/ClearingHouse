package http

import (
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/gin-gonic/gin"
)

func MapNamespaceRoutes(namespaceGroup *gin.RouterGroup, namespaceHandler interfaces.NamespaceHandler) {
	namespaceGroup.POST("/", namespaceHandler.CreateNamespace())
	namespaceGroup.GET("/", namespaceHandler.GetAllNamespaces())
}

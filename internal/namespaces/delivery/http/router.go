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
	namespaceGroup.GET("/all", namespaceHandler.GetAllUserNamespaces())
	namespaceGroup.GET("/:id", namespaceHandler.GetNamespace())
	// namespaceGroup.GET("/:id/quota", namespaceHandler.GetNamespaceQuota())
	// namespaceGroup.GET("/:id/usage", namespaceHandler.GetNamespaceUsage())
	namespaceGroup.POST("/members", namespaceHandler.AddMembers())
}

// get all user's available  namespace
// get namespace detail
// get quota of ns
// get current usage of ns

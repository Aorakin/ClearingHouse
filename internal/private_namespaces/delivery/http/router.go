package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/private_namespaces/interfaces"
	"github.com/gin-gonic/gin"
)

func MapPrivateNamespaceRoutes(privateNamespaceGroup *gin.RouterGroup, privateNamespaceHandler interfaces.PrivateNamespaceHandler) {
	privateNamespaceGroup.Use(middleware.AuthMiddleware())
	privateNamespaceGroup.GET("/", privateNamespaceHandler.GetPrivateNamespace())
	privateNamespaceGroup.POST("/", privateNamespaceHandler.CreatePrivateNamespace())
}

// create priv ns
// create ns quota
// modify ns quota
// user get priv ns
// get usage **may use existing ns usage**
// create ticket **may use existing ticket creation**

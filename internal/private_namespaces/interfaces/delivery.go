package interfaces

import "github.com/gin-gonic/gin"

type PrivateNamespaceHandler interface {
	GetPrivateNamespace() gin.HandlerFunc
	CreatePrivateNamespace() gin.HandlerFunc
}

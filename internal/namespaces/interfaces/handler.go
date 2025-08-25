package interfaces

import (
	"github.com/gin-gonic/gin"
)

type NamespaceHandler interface {
	CreateNamespace() gin.HandlerFunc
	GetAllNamespaces() gin.HandlerFunc
}

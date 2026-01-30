package interfaces

import (
	"github.com/gin-gonic/gin"
)

type NamespaceHandler interface {
	CreateNamespace() gin.HandlerFunc
	GetAllNamespaces() gin.HandlerFunc
	AddMembers() gin.HandlerFunc
	RemoveMembers() gin.HandlerFunc
	UpdateNamespace() gin.HandlerFunc
	DeleteNamespace() gin.HandlerFunc

	GetAllUserNamespaces() gin.HandlerFunc
	GetNamespace() gin.HandlerFunc
	GetNamespaceUsage() gin.HandlerFunc

	GetAllPrivateNamespaces() gin.HandlerFunc
	GetNamespacesByProjectID() gin.HandlerFunc
}

package interfaces

import "github.com/gin-gonic/gin"

type AdminHandler interface {
	AssignSuperAdmin() gin.HandlerFunc
	RevokeSuperAdmin() gin.HandlerFunc
}

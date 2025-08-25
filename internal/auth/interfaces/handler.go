package interfaces

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	GoogleCallback() gin.HandlerFunc
	GoogleLogin() gin.HandlerFunc
	GetMe() gin.HandlerFunc
}

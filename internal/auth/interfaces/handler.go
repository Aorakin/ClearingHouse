package interfaces

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	GoogleCallback() gin.HandlerFunc
	GoogleLogin() gin.HandlerFunc
	Logout() gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
	GetMe() gin.HandlerFunc
}

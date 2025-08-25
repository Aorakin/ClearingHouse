package interfaces

import "github.com/gin-gonic/gin"

type UsersHandlers interface {
	LoginWithGoogle() gin.HandlerFunc
	Callback() gin.HandlerFunc
	Logout() gin.HandlerFunc
	TestSession() gin.HandlerFunc
}

package interfaces

import "github.com/gin-gonic/gin"

type UsersHandlers interface {
	Login() gin.HandlerFunc
	Callback() gin.HandlerFunc
}

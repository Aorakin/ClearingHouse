package http

import (
	"github.com/ClearingHouse/internal/users/interfaces"
	"github.com/gin-gonic/gin"
)

func MapUsersRoutes(usersGroup *gin.RouterGroup, usersHandler interfaces.UsersHandlers) {
	usersGroup.GET("/auth/google", usersHandler.Login())
	usersGroup.GET("/auth/callback/google", usersHandler.Callback())
}

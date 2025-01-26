package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/users/interfaces"
	"github.com/gin-gonic/gin"
)

func MapUsersRoutes(usersGroup *gin.RouterGroup, usersHandler interfaces.UsersHandlers) {
	usersGroup.GET("/auth/google", usersHandler.LoginWithGoogle())
	usersGroup.GET("/auth/callback/google", usersHandler.Callback())
	usersGroup.GET("/register", usersHandler.Register())
	usersGroup.POST("/login", usersHandler.Login())
	usersGroup.POST("/logout", usersHandler.Logout())
	usersGroup.GET("/testsession", middleware.AuthMiddleware(), usersHandler.TestSession())
}

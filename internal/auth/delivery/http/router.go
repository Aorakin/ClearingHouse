package http

import (
	"github.com/ClearingHouse/internal/auth/interfaces"
	"github.com/ClearingHouse/internal/middleware"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, authHandler interfaces.AuthHandler) {
	authGroup.GET("/callback/google", authHandler.GoogleCallback())
	authGroup.GET("/login/google", authHandler.GoogleLogin())
	authGroup.GET("/refresh-token", authHandler.RefreshToken())
	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/me", authHandler.GetMe())
}

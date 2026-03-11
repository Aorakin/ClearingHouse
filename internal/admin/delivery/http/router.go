package http

import (
	"github.com/ClearingHouse/internal/admin/interfaces"
	"github.com/ClearingHouse/internal/middleware"
	"github.com/gin-gonic/gin"
)

func MapAdminRoutes(adminGroup *gin.RouterGroup, adminHandler interfaces.AdminHandler) {
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.SuperAdminMiddleware())
	adminGroup.POST("/super-admins", adminHandler.AssignSuperAdmin())
	adminGroup.POST("/revoke-super-admins", adminHandler.RevokeSuperAdmin())
}

package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/gin-gonic/gin"
)

func MapQuotaRoutes(quotaGroup *gin.RouterGroup, quotaHandler interfaces.QuotaHandler) {
	quotaGroup.Use(middleware.AuthMiddleware())
	quotaGroup.POST("/organization", quotaHandler.CreateOrganizationQuotaGroup())
	quotaGroup.GET("/organization", quotaHandler.FindOrganizationQuotaGroup())
	quotaGroup.POST("/project", quotaHandler.CreateProjectQuotaGroup())
	quotaGroup.GET("/project/:id", quotaHandler.FindProjectQuotaGroup())
	quotaGroup.POST("/namespace", quotaHandler.CreateNamespaceQuotaGroup())
	quotaGroup.POST("/namespace/assign", quotaHandler.AssignQuotaToNamespace())
}

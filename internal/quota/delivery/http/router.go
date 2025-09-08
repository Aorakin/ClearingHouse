package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/gin-gonic/gin"
)

func MapQuotaRoutes(quotaGroup *gin.RouterGroup, quotaHandler interfaces.QuotaHandler) {
	// 	quotaGroup.GET("/namespace/:id", quotaHandler.FindNamespaceQuotaGroup())
	quotaGroup.Use(middleware.AuthMiddleware())
	quotaGroup.POST("/organization", quotaHandler.CreateOrganizationQuota())
	quotaGroup.GET("/organization", quotaHandler.GetOrganizationQuota())
	quotaGroup.POST("/project", quotaHandler.CreateProjectQuota())
	quotaGroup.POST("/project/own", quotaHandler.CreateOwnedProjectQuota())
	quotaGroup.GET("/project/:id", quotaHandler.GetProjectQuota())
	quotaGroup.POST("/namespace", quotaHandler.CreateNamespaceQuota())
	quotaGroup.GET("/namespace/:id", quotaHandler.GetNamespaceQuota())
	quotaGroup.POST("/namespace/assign", quotaHandler.AssignQuotaToNamespace())
	// quotaGroup.POST("/namespace/assign", quotaHandler.AssignQuotaToNamespace())
}

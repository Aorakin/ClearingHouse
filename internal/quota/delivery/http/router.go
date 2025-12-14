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
	quotaGroup.POST("/internal-project/", quotaHandler.CreateOwnedProjectQuota())
	quotaGroup.GET("/project/:project_id", quotaHandler.GetProjectQuota())
	quotaGroup.GET("/project/:project_id/namespaces", quotaHandler.GetNamespaceQuotaInProject())
	quotaGroup.POST("/namespace", quotaHandler.CreateNamespaceQuota())
	quotaGroup.GET("/namespace/:namespace_id", quotaHandler.GetNamespaceQuota())
	quotaGroup.POST("/namespace/assign", quotaHandler.AssignQuotaToNamespace())
	quotaGroup.GET("/:quota_id/usage/:namespace_id", quotaHandler.GetUsage())
	// quotaGroup.POST("/namespace/assign", quotaHandler.AssignQuotaToNamespace())
}

package http

import (
	"github.com/ClearingHouse/internal/middleware"
	"github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/gin-gonic/gin"
)

func MapQuotaRoutes(quotaGroup *gin.RouterGroup, quotaHandler interfaces.QuotaHandler) {
	quotaGroup.Use(middleware.AuthMiddleware())
	quotaGroup.POST("/organization", quotaHandler.CreateOrganizationQuota())
	quotaGroup.GET("/organization", quotaHandler.GetOrganizationQuota())

	quotaGroup.POST("/project", quotaHandler.CreateProjectQuota())
	quotaGroup.GET("/project/:project_id", quotaHandler.GetProjectQuotas())
	quotaGroup.GET("/project/:project_id/namespaces", quotaHandler.GetNamespaceQuotaInProject())
	quotaGroup.POST("/project/internal", quotaHandler.CreateInternalProjectQuota())

	quotaGroup.POST("/namespace", quotaHandler.CreateNamespaceQuota())
	quotaGroup.GET("/namespace/:namespace_id", quotaHandler.GetNamespaceQuota())
	quotaGroup.POST("/namespace/template", quotaHandler.CreateNamespaceQuotaTemplate())
	quotaGroup.GET("/namespace/template/:quota_template_id", quotaHandler.GetNamespaceQuotaTemplate())
	quotaGroup.POST("/namespace/template/assign", quotaHandler.AssignQuotaTemplateToNamespace())

	quotaGroup.GET("/:quota_id/usage/:namespace_id", quotaHandler.GetUsage())
}

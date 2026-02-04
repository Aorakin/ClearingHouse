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
	quotaGroup.GET("/organization/:org_id", quotaHandler.GetOrganizationQuotasByOrgID())
	quotaGroup.DELETE("/organization/:quota_id", quotaHandler.DeleteOrganizationQuota())

	quotaGroup.POST("/project", quotaHandler.CreateProjectQuota())
	quotaGroup.GET("/project/:project_id", quotaHandler.GetProjectQuotas())
	quotaGroup.GET("/project/:project_id/namespaces", quotaHandler.GetNamespaceQuotaInProject())
	quotaGroup.POST("/project/internal", quotaHandler.CreateInternalProjectQuota())
	quotaGroup.DELETE("/project/:quota_id", quotaHandler.DeleteProjectQuota())

	quotaGroup.POST("/namespace", quotaHandler.CreateNamespaceQuota())
	quotaGroup.GET("/namespace/:namespace_id", quotaHandler.GetNamespaceQuota())
	quotaGroup.POST("/namespace/template", quotaHandler.CreateNamespaceQuotaTemplate())
	quotaGroup.GET("/namespace/template/:quota_template_id", quotaHandler.GetNamespaceQuotaTemplate())
	quotaGroup.GET("/namespace/template/project/:project_id", quotaHandler.GetNamespaceQuotaTemplatesByProjectID())
	quotaGroup.PUT("/namespace/template/:quota_template_id", quotaHandler.UpdateNamespaceQuotaTemplate())
	quotaGroup.POST("/namespace/template/assign", quotaHandler.AssignQuotaTemplateToNamespace())
	quotaGroup.DELETE("/namespace/template/unassign/:namespace_id", quotaHandler.UnassignQuotaTemplateFromNamespace())
	quotaGroup.DELETE("/namespace/template/:quota_template_id", quotaHandler.DeleteNamespaceQuotaTemplate())
	quotaGroup.DELETE("/namespace/:quota_id", quotaHandler.DeleteNamespaceQuota())

	quotaGroup.GET("/:quota_id/usage/:namespace_id", quotaHandler.GetUsage())
}

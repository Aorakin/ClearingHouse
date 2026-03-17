package interfaces

import "github.com/gin-gonic/gin"

type QuotaHandler interface {
	CreateOrganizationQuota() gin.HandlerFunc
	GetOrganizationQuota() gin.HandlerFunc
	GetOrganizationQuotasByOrgID() gin.HandlerFunc
	UpdateOrganizationQuota() gin.HandlerFunc
	DeleteOrganizationQuota() gin.HandlerFunc

	CreateProjectQuota() gin.HandlerFunc
	CreateInternalProjectQuota() gin.HandlerFunc
	GetProjectQuotas() gin.HandlerFunc
	GetProjectQuotaTotal() gin.HandlerFunc
	GetNamespaceQuotaInProject() gin.HandlerFunc
	UpdateProjectQuota() gin.HandlerFunc
	UpdateInternalProjectQuota() gin.HandlerFunc
	DeleteProjectQuota() gin.HandlerFunc

	CreateNamespaceQuota() gin.HandlerFunc
	GetNamespaceQuota() gin.HandlerFunc
	UpdateNamespaceQuota() gin.HandlerFunc
	CreateNamespaceQuotaTemplate() gin.HandlerFunc
	GetNamespaceQuotaTemplate() gin.HandlerFunc
	GetNamespaceQuotaTemplatesByProjectID() gin.HandlerFunc
	UpdateNamespaceQuotaTemplate() gin.HandlerFunc
	AssignQuotaTemplateToNamespace() gin.HandlerFunc
	UnassignQuotaTemplateFromNamespace() gin.HandlerFunc
	DeleteNamespaceQuota() gin.HandlerFunc
	DeleteNamespaceQuotaTemplate() gin.HandlerFunc

	GetUsage() gin.HandlerFunc

	// FindOrganizationQuotaGroup() gin.HandlerFunc
	// CreateOrganizationQuotaGroup() gin.HandlerFunc
	// FindProjectQuotaGroup() gin.HandlerFunc
	// CreateProjectQuotaGroup() gin.HandlerFunc
	// FindNamespaceQuotaGroup() gin.HandlerFunc
	// CreateNamespaceQuotaGroup() gin.HandlerFunc
	// AssignQuotaToNamespace() gin.HandlerFunc

	// GetOrganizationQuota() gin.HandlerFunc
	// CreateOrganizationQuota() gin.HandlerFunc
	// GetProjectQuota() gin.HandlerFunc
	// CreateProjectQuota() gin.HandlerFunc
	// GetNamespaceQuota() gin.HandlerFunc
	// CreateNamespaceQuota() gin.HandlerFunc
	// AssignNamespaceQuota() gin.HandlerFunc
}

package interfaces

import "github.com/gin-gonic/gin"

type QuotaHandler interface {
	CreateOrganizationQuota() gin.HandlerFunc
	GetOrganizationQuota() gin.HandlerFunc

	CreateProjectQuota() gin.HandlerFunc
	GetProjectQuota() gin.HandlerFunc

	CreateNamespaceQuota() gin.HandlerFunc
	GetNamespaceQuota() gin.HandlerFunc
	AssignQuotaToNamespace() gin.HandlerFunc

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

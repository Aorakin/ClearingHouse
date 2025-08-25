package interfaces

import (
	"github.com/gin-gonic/gin"
)

type QuotaHandler interface {
	// CreateOrganizationQuota() gin.HandlerFunc
	FindOrganizationQuotaGroup() gin.HandlerFunc
	CreateOrganizationQuotaGroup() gin.HandlerFunc
	FindProjectQuotaGroup() gin.HandlerFunc
	CreateProjectQuotaGroup() gin.HandlerFunc
	CreateNamespaceQuotaGroup() gin.HandlerFunc
	AssignQuotaToNamespace() gin.HandlerFunc
}

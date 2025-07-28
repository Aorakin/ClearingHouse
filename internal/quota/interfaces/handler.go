package interfaces

import (
	"github.com/gin-gonic/gin"
)

type QuotaHandler interface {
	// CreateOrganizationQuota() gin.HandlerFunc
	FindOrganizationQuotaGroup() gin.HandlerFunc
	CreateOrganizationQuotaGroup() gin.HandlerFunc
}

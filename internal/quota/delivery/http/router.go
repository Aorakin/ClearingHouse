package http

import (
	"github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/gin-gonic/gin"
)

func MapQuotaRoutes(quotaGroup *gin.RouterGroup, quotaHandler interfaces.QuotaHandler) {
	quotaGroup.POST("/organization", quotaHandler.CreateOrganizationQuotaGroup())
	quotaGroup.GET("/organization", quotaHandler.FindOrganizationQuotaGroup())
}

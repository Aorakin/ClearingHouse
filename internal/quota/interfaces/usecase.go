package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/google/uuid"
)

type QuotaUsecase interface {
	// CreateOrganizationQuota(request *dtos.CreateOrganizationQuotaRequest) (*models.OrganizationQuotaGroup, error)
	// CreateProjectQuota(request *dtos.CreateProjectQuotaRequest) (*models.ProjectQuotaGroup, error)
	// CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest) (*models.NamespaceQuotaGroup, error)
	FindOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuotaGroup, error)
	CreateOrganizationQuotaGroup(request *dtos.CreateOrganizationQuotaRequest) (*models.OrganizationQuotaGroup, error)
}

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
	FindProjectQuotaGroup(projectID uuid.UUID) ([]models.ProjectQuotaGroup, error)
	CreateProjectQuotaGroup(request *dtos.CreateProjectQuotaRequest) (*models.ProjectQuotaGroup, error)
	CreateNamespaceQuotaGroup(request *dtos.CreateNamespaceQuotaRequest) (*models.NamespaceQuotaGroup, error)
	AssignQuotaToNamespace(request *dtos.AssignQuotaToNamespaceRequest) error
}

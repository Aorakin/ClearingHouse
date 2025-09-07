package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/google/uuid"
)

type QuotaUsecase interface {
	CreateOrganizationQuota(request *dtos.CreateOrganizationQuotaRequest, userID uuid.UUID) (*models.OrganizationQuota, error)
	GetOrganizationQuota(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuota, error)

	CreateProjectQuota(request *dtos.CreateProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error)
	GetProjectQuota(projectID uuid.UUID) ([]models.ProjectQuota, error)

	CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error)
	GetNamespaceQuota(namespaceID uuid.UUID) ([]models.NamespaceQuota, error)

	// GetProjectQuota(projectID uuid.UUID) ([]models.ProjectQuota, error)
	// CreateProjectQuota(request *dtos.CreateProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error)
	// GetNamespaceQuota(namespaceID uuid.UUID) ([]models.NamespaceQuota, error)
	// CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error)
	// AssignNamespaceQuota() gin.HandlerFunc

	// isOrgAdmin(orgID uuid.UUID, userID uuid.UUID) error
}

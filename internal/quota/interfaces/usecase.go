package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/google/uuid"
)

type QuotaUsecase interface {
	CreateOrganizationQuota(request *dtos.CreateOrganizationQuotaRequest, userID uuid.UUID) (*models.OrganizationQuota, error)
	GetOrganizationQuota(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuota, error)
	GetOrganizationQuotasByOrgID(orgID uuid.UUID) ([]models.OrganizationQuota, error)
	DeleteOrganizationQuota(quotaID uuid.UUID, userID uuid.UUID) error

	GetNamespaceQuotaInProject(userID uuid.UUID, projectID uuid.UUID) ([]dtos.NamespaceQuotaResponse, error)
	CreateProjectQuota(request *dtos.CreateProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error)
	GetProjectQuotas(projectID uuid.UUID) ([]models.ProjectQuota, error)
	CreateInternalProjectQuota(request *dtos.CreateInternalProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error)
	DeleteProjectQuota(quotaID uuid.UUID, userID uuid.UUID) error

	CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error)
	GetNamespaceQuota(namespaceID uuid.UUID) ([]dtos.NamespaceQuotaResponse, error)
	UpdateNamespaceQuota(quotaID uuid.UUID, request *dtos.UpdateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error)
	CreateNamespaceQuotaTemplate(request *dtos.CreateNamespaceQuotaTemplateRequest, userID uuid.UUID) (*models.NamespaceQuotaTemplate, error)
	GetNamespaceQuotaTemplate(quotaTemplateID uuid.UUID) (*models.NamespaceQuotaTemplate, error)
	GetNamespaceQuotaTemplatesByProjectID(projectID uuid.UUID, userID uuid.UUID) ([]models.NamespaceQuotaTemplate, error)
	UpdateNamespaceQuotaTemplate(quotaTemplateID uuid.UUID, request *dtos.UpdateNamespaceQuotaTemplateRequest, userID uuid.UUID) (*models.NamespaceQuotaTemplate, error)
	AssignQuotaTemplateToNamespace(request *dtos.AssignQuotaToNamespaceRequest, userID uuid.UUID) error
	UnassignQuotaTemplateFromNamespace(namespaceID uuid.UUID, userID uuid.UUID) error
	DeleteNamespaceQuota(quotaID uuid.UUID, userID uuid.UUID) error
	DeleteNamespaceQuotaTemplate(templateID uuid.UUID, userID uuid.UUID) error

	GetUsage(quotaID uuid.UUID, namespaceID uuid.UUID, userID uuid.UUID) (interface{}, error)

	// GetProjectQuota(projectID uuid.UUID) ([]models.ProjectQuota, error)
	// CreateProjectQuota(request *dtos.CreateProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error)
	// GetNamespaceQuota(namespaceID uuid.UUID) ([]models.NamespaceQuota, error)
	// CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error)
	// AssignNamespaceQuota() gin.HandlerFunc

	// isOrgAdmin(orgID uuid.UUID, userID uuid.UUID) error
}

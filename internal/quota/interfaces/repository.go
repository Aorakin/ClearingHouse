package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/google/uuid"
)

type QuotaRepository interface {
	IsOrgQuotaExist(fromOrgID uuid.UUID, toOrgID uuid.UUID, nodeID uuid.UUID) (bool, error)
	CreateOrgQuota(quota *models.OrganizationQuota) error
	GetOrganizationByRelationship(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuota, error)
	GetOrganizationQuotasByOrgID(orgID uuid.UUID) ([]models.OrganizationQuota, error)
	GetOrgQuotaByID(id uuid.UUID) (*models.OrganizationQuota, error)
	GetOrgUsage(quotaID uuid.UUID, resourceID uuid.UUID) (uint, error)
	GetOrgQuotaQuantity(quotaID uuid.UUID, resourceID uuid.UUID) (uint, error)
	DeleteOrganizationQuotasByOrgID(orgID uuid.UUID) error
	DeleteOrganizationQuota(quotaID uuid.UUID) error
	HasActiveQuotaUsage(orgID uuid.UUID) (bool, error)
	HasProjectQuotasByOrgQuotaID(orgQuotaID uuid.UUID) (bool, error)
	HasActiveUsageByOrgQuotaID(orgQuotaID uuid.UUID) (bool, error)

	IsProjectQuotaExist(projectID uuid.UUID, nodeID uuid.UUID) (bool, error)
	CreateProjectQuota(quota *models.ProjectQuota) error
	GetProjectQuotaByProjectID(projectID uuid.UUID) ([]models.ProjectQuota, error)
	GetProjectQuotaByID(id uuid.UUID) (*models.ProjectQuota, error)
	DeleteProjectQuota(quotaID uuid.UUID) error
	HasNamespaceQuotasByProjectQuotaID(projectQuotaID uuid.UUID) (bool, error)

	IsNamespaceQuotaExists(namespaceID uuid.UUID, nodeID uuid.UUID) (bool, error)
	CreateNamespaceQuota(quota *models.NamespaceQuota) error
	GetNamespaceQuotaByNamespaceID(namespaceID uuid.UUID) ([]models.NamespaceQuota, error)
	GetNamespaceQuotaByID(id uuid.UUID) (*models.NamespaceQuota, error)
	DeleteNamespaceQuota(quotaID uuid.UUID) error
	HasQuotaTemplatesByNamespaceQuotaID(namespaceQuotaID uuid.UUID) (bool, error)

	CreateResourceProperty(resourceProperty *models.ResourceProperty) error

	CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error

	GetNamespaceUsageByType(namespaceID uuid.UUID, quotaID uuid.UUID) (*dtos.ResourceUsageResponse, error)
	GetNamespaceQuotaByType(namespaceID uuid.UUID) (*dtos.ResourceQuotaResponse, error)
	GetQuotaByType(quotaID uuid.UUID) (*dtos.ResourceQuotaResponse, error)

	GetNamespaceQuotasByProjectID(projectID uuid.UUID) ([]models.NamespaceQuota, error)
	GetNamespaceQuotasByIDs(quotaIDs []uuid.UUID, projectID uuid.UUID) ([]models.NamespaceQuota, error)
	CreateNamespaceQuotaTemplate(template *models.NamespaceQuotaTemplate) error
	AddQuotasToTemplate(templateID uuid.UUID, quotaIDs []uuid.UUID) error
	RemoveQuotasFromTemplate(templateID uuid.UUID, quotaIDs []uuid.UUID) error
	GetNamespaceQuotaTemplateByID(templateID uuid.UUID) (*models.NamespaceQuotaTemplate, error)
	GetNamespaceQuotaTemplatesByProjectID(projectID uuid.UUID) ([]models.NamespaceQuotaTemplate, error)
	AssignQuotaToNamespace(namespaceID uuid.UUID, quotaTemplateID uuid.UUID) error
	IsAssigned(namespaceID uuid.UUID, quotaID uuid.UUID) (bool, error)
	DeleteNamespaceQuotaTemplate(templateID uuid.UUID) error
	HasNamespacesUsingTemplate(templateID uuid.UUID) (bool, error)
}

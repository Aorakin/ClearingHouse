package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type QuotaRepository interface {
	IsOrgQuotaExist(fromOrgID uuid.UUID, toOrgID uuid.UUID, poolID uuid.UUID) (bool, error)
	CreateOrgQuota(quota *models.OrganizationQuota) error
	GetOrganizationByRelationship(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuota, error)
	GetOrgQuotaByID(id uuid.UUID) (*models.OrganizationQuota, error)
	GetOrgUsage(quotaID uuid.UUID, resourceID uuid.UUID) (uint, error)
	GetOrgQuotaQuantity(quotaID uuid.UUID, resourceID uuid.UUID) (uint, error)

	CreateProjectQuota(quota *models.ProjectQuota) error
	GetProjectQuotaByProjectID(projectID uuid.UUID) ([]models.ProjectQuota, error)
	GetProjectQuotaByID(id uuid.UUID) (*models.ProjectQuota, error)

	CreateNamespaceQuota(quota *models.NamespaceQuota) error
	GetNamespaceQuotaByNamespaceID(namespaceID uuid.UUID) ([]models.NamespaceQuota, error)

	CreateResourceProperty(resourceProperty *models.ResourceProperty) error

	CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error

	// FindOrganizationQuotaGroup(fromOrgId uuid.UUID, toOrgId uuid.UUID) ([]models.OrganizationQuotaGroup, error)
	// FindOrganizationQuotaGroupByID(id uuid.UUID) (*models.OrganizationQuotaGroup, error)
	// FindExistingOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID, poolID uuid.UUID) (*models.OrganizationQuotaGroup, error)
	// CreateOrganizationQuotaGroup(quota *models.OrganizationQuotaGroup) error

	// CreateNamespaceQuotaGroup(quota *models.NamespaceQuotaGroup) error
	// CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error
	// CreateResourceProperty(resourceProperty *models.ResourceProperty) error

	// GetOrgUsage(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error)
	// GetOrgQuotaQuantity(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error)
	// GetResourcePropertyByOrg(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (*models.ResourceProperty, error)

	// FindProjectQuotaGroupByID(id uuid.UUID) (*models.ProjectQuotaGroup, error)
	// CreateProjectQuotaGroup(quota *models.ProjectQuotaGroup) error
	// FindProjectQuotaGroupByProjectID(projectID uuid.UUID) ([]models.ProjectQuotaGroup, error)

	// FindNamespaceQuotaGroupByID(id uuid.UUID) (*models.NamespaceQuotaGroup, error)
	// GetNamespaceQuotaQuantity(namespaceQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error)
	// AssignQuotaToNamespace(namespaceID uuid.UUID, quotaGroupID uuid.UUID) error
	// GetResourcePropertyByNamespace(namespaceQuotaGroupID uuid.UUID, resourceID uuid.UUID) (*models.ResourceProperty, error)

	// GetProjectQuotaQuantity(projQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error)
	// GetResourcePropertyByProj(projQuotaGroupID uuid.UUID, resourceID uuid.UUID) (*models.ResourceProperty, error)
}

package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type QuotaRepository interface {
	FindOrganizationQuotaGroup(fromOrgId uuid.UUID, toOrgId uuid.UUID) ([]models.OrganizationQuotaGroup, error)
	FindOrganizationQuotaGroupByID(id uuid.UUID) (*models.OrganizationQuotaGroup, error)
	FindExistingOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID, poolID uuid.UUID) (*models.OrganizationQuotaGroup, error)
	CreateOrganizationQuotaGroup(quota *models.OrganizationQuotaGroup) error

	CreateNamespaceQuotaGroup(quota *models.NamespaceQuotaGroup) error
	CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error
	CreateResourceProperty(resourceProperty *models.ResourceProperty) error

	GetOrgUsage(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error)
	GetOrgQuotaQuantity(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error)
	GetResourcePropertyByOrg(orgQuotaGroupID uuid.UUID, resourceID uuid.UUID) (*models.ResourceProperty, error)

	FindProjectQuotaGroupByID(id uuid.UUID) (*models.ProjectQuotaGroup, error)
	CreateProjectQuotaGroup(quota *models.ProjectQuotaGroup) error
	FindProjectQuotaGroupByProjectID(projectID uuid.UUID) ([]models.ProjectQuotaGroup, error)

	FindNamespaceQuotaGroupByID(id uuid.UUID) (*models.NamespaceQuotaGroup, error)
	AssignQuotaToNamespace(namespaceID uuid.UUID, quotaGroupID uuid.UUID) error

	GetProjectQuotaQuantity(projQuotaGroupID uuid.UUID, resourceID uuid.UUID) (uint, error)
	GetResourcePropertyByProj(projQuotaGroupID uuid.UUID, resourceID uuid.UUID) (*models.ResourceProperty, error)
}

package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type QuotaRepository interface {
	FindOrganizationQuotaGroup(fromOrgId uuid.UUID, toOrgId uuid.UUID) ([]models.OrganizationQuotaGroup, error)
	FindExistingOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID, poolID uuid.UUID) (*models.OrganizationQuotaGroup, error)
	CreateOrganizationQuotaGroup(quota *models.OrganizationQuotaGroup) error
	// CreateProjectQuotaGroup(quota *models.ProjectQuotaGroup) error
	// CreateNamespaceQuotaGroup(quota *models.NamespaceQuotaGroup) error
	CreateResourceQuantity(resourceQuantity *models.ResourceQuantity) error
	CreateResourceProperty(resourceProperty *models.ResourceProperty) error
}

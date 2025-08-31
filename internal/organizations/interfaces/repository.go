package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type OrganizationRepository interface {
	CreateOrganization(org *models.Organization) (*models.Organization, error)
	GetOrganizationByID(id uuid.UUID) (*models.Organization, error)
	UpdateOrganization(org *models.Organization) (*models.Organization, error)
	DeleteOrganization(id uuid.UUID) error
	GetOrganizations() ([]models.Organization, error)
	UpdateMembers(org *models.Organization) error
}

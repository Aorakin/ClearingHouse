package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/google/uuid"
)

type OrganizationUsecase interface {
	CreateOrganization(org *dtos.CreateOrganization) (*models.Organization, error)
	GetOrganizationByID(id uuid.UUID, userID uuid.UUID) (*models.Organization, error)
	UpdateOrganization(id uuid.UUID, org *dtos.UpdateOrganization) (*models.Organization, error)
	DeleteOrganization(id uuid.UUID) error
	GetOrganizations() ([]models.Organization, error)
	AddMembers(request *dtos.AddMembersRequest) (*models.Organization, error)
}

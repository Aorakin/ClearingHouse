package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/google/uuid"
)

type OrganizationUsecase interface {
	GetAllOrganizations() ([]models.Organization, error)
	GetOrganizationByID(orgID uuid.UUID, userID uuid.UUID) (*models.Organization, error)
	CreateOrganization(request *dtos.CreateOrganization, userID uuid.UUID) (*models.Organization, error)
	AddMembers(request *dtos.AddMembersRequest, userID uuid.UUID) (*models.Organization, error)
}

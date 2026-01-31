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
	UpdateOrganization(orgID uuid.UUID, request *dtos.UpdateOrganization, userID uuid.UUID) (*models.Organization, error)
	DeleteOrganization(orgID uuid.UUID, userID uuid.UUID) error
	AddMembers(request *dtos.AddMembersRequest, userID uuid.UUID) (*models.Organization, error)
	RemoveMembers(request *dtos.RemoveMembersRequest, userID uuid.UUID) (*models.Organization, error)
	AddAdmins(request *dtos.AddAdminsRequest, userID uuid.UUID) (*models.Organization, error)
	RemoveAdmins(request *dtos.RemoveAdminsRequest, userID uuid.UUID) (*models.Organization, error)
	GetMembers() ([]models.User, error)
	GetOrganizationMembers(orgID uuid.UUID) ([]models.User, error)
}

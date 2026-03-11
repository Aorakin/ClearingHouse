package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	orgDtos "github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/google/uuid"
)

type OrganizationUsecase interface {
	GetAllOrganizations(userID uuid.UUID, isSuperAdmin bool) ([]orgDtos.OrganizationResponse, error)
	GetOrganizationByID(orgID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) (*orgDtos.OrganizationResponse, error)
	CreateOrganization(request *orgDtos.CreateOrganization, userID uuid.UUID) (*models.Organization, error)
	UpdateOrganization(orgID uuid.UUID, request *orgDtos.UpdateOrganization, userID uuid.UUID) (*models.Organization, error)
	DeleteOrganization(orgID uuid.UUID, userID uuid.UUID) error
	AddMembers(request *orgDtos.AddMembersRequest, userID uuid.UUID) (*models.Organization, error)
	RemoveMembers(request *orgDtos.RemoveMembersRequest, userID uuid.UUID) (*models.Organization, error)
	AddAdmins(request *orgDtos.AddAdminsRequest, userID uuid.UUID, isSuperAdmin bool) (*models.Organization, error)
	RemoveAdmins(request *orgDtos.RemoveAdminsRequest, userID uuid.UUID, isSuperAdmin bool) (*models.Organization, error)
	GetMembers(isSuperAdmin bool) ([]models.User, error)
	GetOrganizationMembers(orgID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) ([]models.User, error)
}

package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apierror "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type OrganizationUsecase struct {
	orgRepo  interfaces.OrganizationRepository
	userRepo userInterfaces.UsersRepository
}

func NewOrganizationUsecase(orgRepo interfaces.OrganizationRepository, userRepo userInterfaces.UsersRepository) interfaces.OrganizationUsecase {
	return &OrganizationUsecase{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

func (u *OrganizationUsecase) GetAllOrganizations() ([]models.Organization, error) {
	orgs, err := u.orgRepo.GetOrganizations()
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}
	return orgs, nil
}

func (u *OrganizationUsecase) GetOrganizationByID(id uuid.UUID, userID uuid.UUID) (*models.Organization, error) {
	organization, err := u.orgRepo.GetOrganizationByID(id)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(organization.Members, userID) {
		return nil, apierror.NewUnauthorizedError("user is not organization member")
	}

	return organization, nil
}

func (u *OrganizationUsecase) CreateOrganization(request *dtos.CreateOrganization, userID uuid.UUID) (*models.Organization, error) {
	creator, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	organization := &models.Organization{
		Name:        request.Name,
		Description: request.Description,
		Admins:      []models.User{*creator},
		Members:     []models.User{*creator},
	}

	org, err := u.orgRepo.CreateOrganization(organization)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	return org, nil
}

func (u *OrganizationUsecase) AddMembers(request *dtos.AddMembersRequest, userID uuid.UUID) (*models.Organization, error) {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(org.Admins, request.Creator) {
		return nil, apierror.NewUnauthorizedError("user is not organization admin")
	}

	existing := make(map[uuid.UUID]struct{})
	for _, m := range org.Members {
		existing[m.ID] = struct{}{}
	}

	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range request.Members {
		if _, found := existing[memberID]; found {
			return nil, apierror.NewConflictError(fmt.Sprintf("user %s is already an organization member", memberID))
		}
		if _, found := seenReq[memberID]; found {
			return nil, apierror.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", memberID))
		}
		seenReq[memberID] = struct{}{}
		if _, err := u.userRepo.GetByID(memberID); err != nil {
			return nil, apierror.NewNotFoundError(fmt.Sprintf("user %s not found", memberID))
		}
	}

	users, err := u.userRepo.GetByIDs(request.Members)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}
	org.Members = append(org.Members, users...)

	if err := u.orgRepo.UpdateMembers(org); err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	return org, nil
}

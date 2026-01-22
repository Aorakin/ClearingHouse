package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apierror "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type OrganizationUsecase struct {
	orgRepo   interfaces.OrganizationRepository
	userRepo  userInterfaces.UsersRepository
	quotaRepo quotaInterfaces.QuotaRepository
}

func NewOrganizationUsecase(orgRepo interfaces.OrganizationRepository, userRepo userInterfaces.UsersRepository, quotaRepo quotaInterfaces.QuotaRepository) interfaces.OrganizationUsecase {
	return &OrganizationUsecase{
		orgRepo:   orgRepo,
		userRepo:  userRepo,
		quotaRepo: quotaRepo,
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
func (u *OrganizationUsecase) UpdateOrganization(orgID uuid.UUID, request *dtos.UpdateOrganization, userID uuid.UUID) (*models.Organization, error) {
	org, err := u.orgRepo.GetOrganizationByID(orgID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(org.Admins, userID) {
		return nil, apierror.NewUnauthorizedError("user is not organization admin")
	}

	org.Name = request.Name
	org.Description = request.Description

	updatedOrg, err := u.orgRepo.UpdateOrganization(org)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	return updatedOrg, nil
}

func (u *OrganizationUsecase) DeleteOrganization(orgID uuid.UUID, userID uuid.UUID) error {
	org, err := u.orgRepo.GetOrganizationByID(orgID)
	if err != nil {
		return apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(org.Admins, userID) {
		return apierror.NewUnauthorizedError("user is not organization admin")
	}

	// Prevent deletion if organization has projects
	if len(org.Projects) > 0 {
		return apierror.NewBadRequestError(fmt.Errorf("cannot delete organization with existing projects. Please delete all projects first"))
	}

	// Prevent deletion if organization has resource pools
	if len(org.ResourcePools) > 0 {
		return apierror.NewBadRequestError(fmt.Errorf("cannot delete organization with existing resource pools. Please delete all resource pools first"))
	}

	// Prevent deletion if organization has given quotas to other organizations
	// These quotas are being used by other organizations and must be revoked first
	if len(org.GivenQuotas) > 0 {
		return apierror.NewBadRequestError(fmt.Errorf("cannot delete organization with given quotas to other organizations. Please revoke all given quotas first"))
	}

	// Delete all quotas received by this organization (ToOrgID = orgID)
	if err := u.quotaRepo.DeleteOrganizationQuotasByOrgID(orgID); err != nil {
		return apierror.NewInternalServerError(err)
	}

	if err := u.orgRepo.DeleteOrganization(orgID); err != nil {
		return apierror.NewInternalServerError(err)
	}

	return nil
}

func (u *OrganizationUsecase) AddMembers(request *dtos.AddMembersRequest, userID uuid.UUID) (*models.Organization, error) {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(org.Admins, userID) {
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

func (u *OrganizationUsecase) RemoveMembers(request *dtos.RemoveMembersRequest, userID uuid.UUID) (*models.Organization, error) {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(org.Admins, userID) {
		return nil, apierror.NewUnauthorizedError("user is not organization admin")
	}

	// Create a map of members to remove for quick lookup
	removeMap := make(map[uuid.UUID]struct{})
	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range request.Members {
		if _, found := seenReq[memberID]; found {
			return nil, apierror.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", memberID))
		}
		seenReq[memberID] = struct{}{}
		removeMap[memberID] = struct{}{}
	}

	// Check if members to remove exist in organization and prevent removing admins
	existing := make(map[uuid.UUID]struct{})
	for _, m := range org.Members {
		existing[m.ID] = struct{}{}
	}

	for _, memberID := range request.Members {
		if _, found := existing[memberID]; !found {
			return nil, apierror.NewNotFoundError(fmt.Sprintf("user %s is not a member of this organization", memberID))
		}
		// Prevent removing admins through member removal
		if helper.ContainsUserID(org.Admins, memberID) {
			return nil, apierror.NewBadRequestError(fmt.Sprintf("user %s is an admin and cannot be removed as a member directly", memberID))
		}
	}

	// Filter out members to remove
	newMembers := []models.User{}
	for _, member := range org.Members {
		if _, shouldRemove := removeMap[member.ID]; !shouldRemove {
			newMembers = append(newMembers, member)
		}
	}

	org.Members = newMembers

	if err := u.orgRepo.UpdateMembers(org); err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	return org, nil
}

func (u *OrganizationUsecase) GetMembers() ([]models.User, error) {
	users, err := u.orgRepo.GetMembers()
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}
	return users, nil
}
func (u *OrganizationUsecase) GetOrganizationMembers(orgID uuid.UUID) ([]models.User, error) {
	users, err := u.orgRepo.GetOrganizationMembers(orgID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}
	return users, nil
}

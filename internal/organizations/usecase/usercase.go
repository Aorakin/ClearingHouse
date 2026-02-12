package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceDtos "github.com/ClearingHouse/internal/namespaces/dtos"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	projectInterfaces "github.com/ClearingHouse/internal/projects/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apierror "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type OrganizationUsecase struct {
	orgRepo       interfaces.OrganizationRepository
	userRepo      userInterfaces.UsersRepository
	quotaRepo     quotaInterfaces.QuotaRepository
	projectRepo   projectInterfaces.ProjectRepository
	namespaceRepo namespaceInterfaces.NamespaceRepository
}

func NewOrganizationUsecase(orgRepo interfaces.OrganizationRepository, userRepo userInterfaces.UsersRepository, quotaRepo quotaInterfaces.QuotaRepository, projectRepo projectInterfaces.ProjectRepository, namespaceRepo namespaceInterfaces.NamespaceRepository) interfaces.OrganizationUsecase {
	return &OrganizationUsecase{
		orgRepo:       orgRepo,
		userRepo:      userRepo,
		quotaRepo:     quotaRepo,
		projectRepo:   projectRepo,
		namespaceRepo: namespaceRepo,
	}
}

func (u *OrganizationUsecase) GetAllOrganizations() ([]models.Organization, error) {
	orgs, err := u.orgRepo.GetOrganizations()
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}
	return orgs, nil
}

func (u *OrganizationUsecase) GetOrganizationByID(id uuid.UUID, userID uuid.UUID) (*dtos.OrganizationResponse, error) {
	organization, err := u.orgRepo.GetOrganizationByID(id)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(organization.Admins, userID) && !helper.ContainsUserID(organization.Members, userID) {
		return nil, apierror.NewUnauthorizedError("user is not organization admin or member")
	}

	// Get all projects for this organization
	projects, err := u.projectRepo.GetProjectsByOrganizationID(id)
	if err != nil {
		return nil, apierror.NewInternalServerError(fmt.Errorf("failed to get projects: %w", err))
	}

	// Aggregate resources by resource type across all projects
	typeAgg := make(map[uuid.UUID]namespaceDtos.ResourceQuota)

	// Loop through each project
	for _, project := range projects {
		// Get all namespaces for this project with preloaded quota templates
		namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(project.ID)
		if err != nil {
			return nil, apierror.NewInternalServerError(fmt.Errorf("failed to get namespaces for project %s: %w", project.ID, err))
		}

		// Loop through namespaces
		for _, namespace := range namespaces {
			// Skip namespaces without quota template
			if namespace.QuotaTemplate == nil {
				continue
			}

			// Iterate through all quotas in the template
			for _, quota := range namespace.QuotaTemplate.Quotas {
				// Iterate through all resources in the quota
				for _, resource := range quota.Resources {
					rt := resource.ResourceProp.Resource.ResourceType
					rtID := rt.ID

					// Initialize if not exists
					if _, ok := typeAgg[rtID]; !ok {
						typeAgg[rtID] = namespaceDtos.ResourceQuota{
							TypeID: rtID,
							Type:   rt.Name,
							Quota:  0,
						}
					}

					// Add to existing quota
					tmp := typeAgg[rtID]
					tmp.Quota += float64(resource.Quantity)
					typeAgg[rtID] = tmp
				}
			}
		}
	}

	// Convert map to slice
	var resourceQuotas []namespaceDtos.ResourceQuota
	for _, v := range typeAgg {
		resourceQuotas = append(resourceQuotas, v)
	}

	// Build response
	response := &dtos.OrganizationResponse{
		ID:             organization.ID,
		CreatedAt:      organization.CreatedAt,
		UpdatedAt:      organization.UpdatedAt,
		Name:           organization.Name,
		Description:    organization.Description,
		Members:        organization.Members,
		Admins:         organization.Admins,
		ResourceQuotas: resourceQuotas,
	}

	return response, nil
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
		Members:     []models.User{},
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

	// Check if organization has any active quota usage
	hasUsage, err := u.quotaRepo.HasActiveQuotaUsage(orgID)
	if err != nil {
		return apierror.NewInternalServerError(err)
	}

	if hasUsage {
		return apierror.NewBadRequestError(fmt.Errorf("cannot delete organization with active quota usage. Please release all resources first"))
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

func (u *OrganizationUsecase) AddAdmins(request *dtos.AddAdminsRequest, userID uuid.UUID) (*models.Organization, error) {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(org.Admins, userID) {
		return nil, apierror.NewUnauthorizedError("user is not organization admin")
	}

	existing := make(map[uuid.UUID]struct{})
	for _, a := range org.Admins {
		existing[a.ID] = struct{}{}
	}

	seenReq := make(map[uuid.UUID]struct{})
	for _, adminID := range request.Admins {
		if _, found := existing[adminID]; found {
			return nil, apierror.NewConflictError(fmt.Sprintf("user %s is already an organization admin", adminID))
		}
		if _, found := seenReq[adminID]; found {
			return nil, apierror.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", adminID))
		}
		seenReq[adminID] = struct{}{}
		if _, err := u.userRepo.GetByID(adminID); err != nil {
			return nil, apierror.NewNotFoundError(fmt.Sprintf("user %s not found", adminID))
		}
	}

	users, err := u.userRepo.GetByIDs(request.Admins)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}
	org.Admins = append(org.Admins, users...)

	if err := u.orgRepo.UpdateAdmins(org); err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	return org, nil
}

func (u *OrganizationUsecase) RemoveAdmins(request *dtos.RemoveAdminsRequest, userID uuid.UUID) (*models.Organization, error) {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(org.Admins, userID) {
		return nil, apierror.NewUnauthorizedError("user is not organization admin")
	}

	// Create a map of admins to remove for quick lookup
	removeMap := make(map[uuid.UUID]struct{})
	seenReq := make(map[uuid.UUID]struct{})
	for _, adminID := range request.Admins {
		if _, found := seenReq[adminID]; found {
			return nil, apierror.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", adminID))
		}
		seenReq[adminID] = struct{}{}
		removeMap[adminID] = struct{}{}
	}

	// Check if admins to remove exist in organization
	existing := make(map[uuid.UUID]struct{})
	for _, a := range org.Admins {
		existing[a.ID] = struct{}{}
	}

	for _, adminID := range request.Admins {
		if _, found := existing[adminID]; !found {
			return nil, apierror.NewNotFoundError(fmt.Sprintf("user %s is not an admin of this organization", adminID))
		}
	}

	// Filter out admins to remove
	newAdmins := []models.User{}
	for _, admin := range org.Admins {
		if _, shouldRemove := removeMap[admin.ID]; !shouldRemove {
			newAdmins = append(newAdmins, admin)
		}
	}

	// Prevent removing all admins
	if len(newAdmins) == 0 {
		return nil, apierror.NewBadRequestError("cannot remove all admins from organization. At least one admin must remain")
	}

	org.Admins = newAdmins

	if err := u.orgRepo.UpdateAdmins(org); err != nil {
		return nil, apierror.NewInternalServerError(err)
	}

	return org, nil
}

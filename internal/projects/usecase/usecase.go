package usecase

import (
	"fmt"
	"log"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceDtos "github.com/ClearingHouse/internal/namespaces/dtos"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/ClearingHouse/internal/projects/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type ProjectUsecase struct {
	projRepo      interfaces.ProjectRepository
	orgRepo       orgInterfaces.OrganizationRepository
	userRepo      userInterfaces.UsersRepository
	namespaceRepo namespaceInterfaces.NamespaceRepository
}

func NewProjectUsecase(projRepo interfaces.ProjectRepository, orgRepo orgInterfaces.OrganizationRepository, userRepo userInterfaces.UsersRepository, namespaceRepo namespaceInterfaces.NamespaceRepository) interfaces.ProjectUsecase {
	return &ProjectUsecase{
		projRepo:      projRepo,
		orgRepo:       orgRepo,
		userRepo:      userRepo,
		namespaceRepo: namespaceRepo,
	}
}

// isProjAdminOrOrgAdmin checks if user is a project admin or an org admin of the project's organization.
func (u *ProjectUsecase) isProjAdminOrOrgAdmin(project *models.Project, userID uuid.UUID) bool {
	if helper.ContainsUserID(project.Admins, userID) {
		return true
	}
	org, err := u.orgRepo.GetOrganizationByID(project.OrganizationID)
	if err != nil {
		return false
	}
	return helper.ContainsUserID(org.Admins, userID)
}

func (u *ProjectUsecase) CreateProject(request *dtos.CreateProjectRequest, userID uuid.UUID) error {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(org.Admins, userID) {
		return apiError.NewUnauthorizedError("user is not org admin")
	}

	admin, err := u.userRepo.GetByID(userID)
	if err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	project := &models.Project{
		Name:           request.Name,
		Description:    request.Description,
		OrganizationID: request.OrganizationID,
		Admins:         []models.User{*admin},
	}

	err = u.projRepo.CreateProject(project)
	if err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	return nil
}

func (u *ProjectUsecase) GetAllProjects(userID uuid.UUID, isSuperAdmin bool) ([]dtos.ProjectResponse, error) {
	var projects []models.Project
	var err error

	if isSuperAdmin {
		projects, err = u.projRepo.GetAllProjects()
	} else {
		projects, err = u.projRepo.GetProjectsByUserAdminScope(userID)
	}
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !isSuperAdmin && len(projects) == 0 {
		return nil, apiError.NewForbiddenError("admin access required")
	}

	var projectResponses []dtos.ProjectResponse
	for _, project := range projects {
		// Get all namespaces for this project
		namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(project.ID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespaces: %w", err).Error())
		}

		// Aggregate resources by resource type
		typeAgg := make(map[uuid.UUID]namespaceDtos.ResourceQuota)

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

		// Convert map to slice
		var resourceQuotas []namespaceDtos.ResourceQuota
		for _, v := range typeAgg {
			resourceQuotas = append(resourceQuotas, v)
		}

		// Build response
		projectResponse := dtos.ProjectResponse{
			ID:             project.ID,
			CreatedAt:      project.CreatedAt,
			UpdatedAt:      project.UpdatedAt,
			Name:           project.Name,
			Description:    project.Description,
			OrganizationID: project.OrganizationID,
			Members:        project.Members,
			Admins:         project.Admins,
			ResourceQuotas: resourceQuotas,
		}
		projectResponses = append(projectResponses, projectResponse)
	}

	return projectResponses, nil
}

func (u *ProjectUsecase) AddMembers(request *dtos.AddMembersRequest, userID uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	org, err := u.orgRepo.GetOrganizationByID(project.OrganizationID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !u.isProjAdminOrOrgAdmin(project, userID) {
		return nil, apiError.NewForbiddenError("user is not a project admin or organization admin")
	}

	existing := make(map[uuid.UUID]struct{})
	for _, m := range project.Members {
		existing[m.ID] = struct{}{}
	}

	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range request.Members {
		// must be org member
		if !helper.ContainsUserID(org.Members, memberID) {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("user %s is not a member of the organization", memberID))
		}
		if _, found := existing[memberID]; found {
			return nil, apiError.NewConflictError(fmt.Sprintf("user %s is already a project member", memberID))
		}
		if _, found := seenReq[memberID]; found {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", memberID))
		}
		seenReq[memberID] = struct{}{}
		if _, err := u.userRepo.GetByID(memberID); err != nil {
			return nil, apiError.NewNotFoundError(fmt.Sprintf("user %s not found", memberID))
		}
	}

	users, err := u.userRepo.GetByIDs(request.Members)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	project.Members = append(project.Members, users...)

	if err := u.projRepo.UpdateMembers(project); err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return project, nil
}

func (u *ProjectUsecase) RemoveMembers(request *dtos.RemoveMembersRequest, userID uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !u.isProjAdminOrOrgAdmin(project, userID) {
		return nil, apiError.NewForbiddenError("user is not a project admin or organization admin")
	}

	// Create a map of members to remove for quick lookup
	removeMap := make(map[uuid.UUID]struct{})
	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range request.Members {
		if _, found := seenReq[memberID]; found {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", memberID))
		}
		seenReq[memberID] = struct{}{}
		removeMap[memberID] = struct{}{}
	}

	// Check if members to remove exist in project and prevent removing admins
	existing := make(map[uuid.UUID]struct{})
	for _, m := range project.Members {
		existing[m.ID] = struct{}{}
	}

	for _, memberID := range request.Members {
		if _, found := existing[memberID]; !found {
			return nil, apiError.NewNotFoundError(fmt.Sprintf("user %s is not a member of this project", memberID))
		}
	}

	// Filter out members to remove
	newMembers := []models.User{}
	for _, member := range project.Members {
		if _, shouldRemove := removeMap[member.ID]; !shouldRemove {
			newMembers = append(newMembers, member)
		}
	}

	project.Members = newMembers

	if err := u.projRepo.UpdateMembers(project); err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	// Keep hierarchy consistent: removing a project member also removes
	// the same user from all namespace memberships in that project.
	namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	for _, namespace := range namespaces {
		updatedNamespaceMembers := make([]models.User, 0, len(namespace.Members))
		for _, member := range namespace.Members {
			if _, shouldRemove := removeMap[member.ID]; !shouldRemove {
				updatedNamespaceMembers = append(updatedNamespaceMembers, member)
			}
		}

		namespace.Members = updatedNamespaceMembers
		if err := u.namespaceRepo.UpdateMembers(&namespace); err != nil {
			return nil, apiError.NewInternalServerError(err.Error())
		}
	}

	return project, nil
}

func (u *ProjectUsecase) GetAllUserProjects(userID uuid.UUID) ([]dtos.ProjectResponse, error) {
	projects, err := u.projRepo.GetAllProjectsByUserID(userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	var projectResponses []dtos.ProjectResponse
	for _, project := range projects {
		// Get all namespaces for this project
		namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(project.ID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespaces: %w", err).Error())
		}

		// Aggregate resources by resource type
		typeAgg := make(map[uuid.UUID]namespaceDtos.ResourceQuota)

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

		// Convert map to slice
		var resourceQuotas []namespaceDtos.ResourceQuota
		for _, v := range typeAgg {
			resourceQuotas = append(resourceQuotas, v)
		}

		// Build response
		projectResponse := dtos.ProjectResponse{
			ID:             project.ID,
			CreatedAt:      project.CreatedAt,
			UpdatedAt:      project.UpdatedAt,
			Name:           project.Name,
			Description:    project.Description,
			OrganizationID: project.OrganizationID,
			Members:        project.Members,
			Admins:         project.Admins,
			ResourceQuotas: resourceQuotas,
		}
		projectResponses = append(projectResponses, projectResponse)
	}

	return projectResponses, nil
}

func (u *ProjectUsecase) GetProjectsByOrganizationID(orgID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) ([]dtos.ProjectResponse, error) {
	var isOrgAdmin bool
	if !isSuperAdmin {
		org, err := u.orgRepo.GetOrganizationByID(orgID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err.Error())
		}

		isOrgAdmin = helper.ContainsUserID(org.Admins, userID)

		if !isOrgAdmin && !helper.ContainsUserID(org.Members, userID) {
			// Allow project admins in this org to view their own projects
			isProjectAdmin, pErr := u.orgRepo.IsProjectAdminInOrg(orgID, userID)
			if pErr != nil || !isProjectAdmin {
				return nil, apiError.NewForbiddenError("access denied")
			}
		}
	}

	projects, err := u.projRepo.GetProjectsByOrganizationID(orgID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	var projectResponses []dtos.ProjectResponse
	for _, project := range projects {
		// If not super admin and not org admin, only show projects where user is admin or member
		if !isSuperAdmin && !isOrgAdmin {
			if !helper.ContainsUserID(project.Admins, userID) && !helper.ContainsUserID(project.Members, userID) {
				continue
			}
		}

		// Get all namespaces for this project
		namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(project.ID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespaces: %w", err).Error())
		}

		// Aggregate resources by resource type
		typeAgg := make(map[uuid.UUID]namespaceDtos.ResourceQuota)

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

		// Convert map to slice
		var resourceQuotas []namespaceDtos.ResourceQuota
		for _, v := range typeAgg {
			resourceQuotas = append(resourceQuotas, v)
		}

		// Build response
		projectResponse := dtos.ProjectResponse{
			ID:             project.ID,
			CreatedAt:      project.CreatedAt,
			UpdatedAt:      project.UpdatedAt,
			Name:           project.Name,
			Description:    project.Description,
			OrganizationID: project.OrganizationID,
			Members:        project.Members,
			Admins:         project.Admins,
			ResourceQuotas: resourceQuotas,
		}
		projectResponses = append(projectResponses, projectResponse)
	}

	return projectResponses, nil
}

func (u *ProjectUsecase) GetProjectMembers(projectID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) ([]models.User, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !isSuperAdmin && !helper.ContainsUserID(project.Admins, userID) && !helper.ContainsUserID(project.Members, userID) {
		// Also allow org admin
		org, orgErr := u.orgRepo.GetOrganizationByID(project.OrganizationID)
		if orgErr != nil || !helper.ContainsUserID(org.Admins, userID) {
			return nil, apiError.NewForbiddenError("access denied")
		}
	}

	members, err := u.projRepo.GetProjectMembers(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return members, nil
}

func (u *ProjectUsecase) GetProjectByID(projectID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) (*dtos.ProjectResponse, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !isSuperAdmin && !helper.ContainsUserID(project.Admins, userID) && !helper.ContainsUserID(project.Members, userID) {
		// Allow org admins to view
		org, err := u.orgRepo.GetOrganizationByID(project.OrganizationID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err.Error())
		}
		if !helper.ContainsUserID(org.Admins, userID) {
			return nil, apiError.NewForbiddenError("access denied")
		}
	}

	// Get all namespaces for this project with preloaded quota templates
	namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespaces: %w", err).Error())
	}

	// Aggregate resources by resource type
	typeAgg := make(map[uuid.UUID]namespaceDtos.ResourceQuota)

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

	// Convert map to slice
	var resourceQuotas []namespaceDtos.ResourceQuota
	for _, v := range typeAgg {
		resourceQuotas = append(resourceQuotas, v)
	}

	// Build response
	response := &dtos.ProjectResponse{
		ID:             project.ID,
		CreatedAt:      project.CreatedAt,
		UpdatedAt:      project.UpdatedAt,
		Name:           project.Name,
		Description:    project.Description,
		OrganizationID: project.OrganizationID,
		Members:        project.Members,
		Admins:         project.Admins,
		ResourceQuotas: resourceQuotas,
	}

	return response, nil
}

func (u *ProjectUsecase) GetProjectUsage(projectID uuid.UUID, userID uuid.UUID) (*dtos.ProjectUsageResponse, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(project.Admins, userID) && !helper.ContainsUserID(project.Members, userID) {
		// Also allow org admin
		org, orgErr := u.orgRepo.GetOrganizationByID(project.OrganizationID)
		if orgErr != nil || !helper.ContainsUserID(org.Admins, userID) {
			log.Println("Unauthorized access attempt by user:", userID)
			return nil, apiError.NewForbiddenError("user is not a project admin, member, or organization admin")
		}
	}

	quotas, err := u.projRepo.GetProjectQuotaByType(projectID, userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	usages, err := u.projRepo.GetProjectUsageByType(projectID, userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	var projectUsage dtos.ProjectUsageResponse

	usageMap := make(map[string]float64)
	for _, u := range usages.ResourceUsages {
		usageMap[u.TypeID] = u.Usage
	}

	for _, q := range quotas.ResourceQuotas {
		uVal := usageMap[q.TypeID]
		projectUsage.Usage = append(projectUsage.Usage, dtos.ProjectUsage{
			TypeID: q.TypeID,
			Type:   q.Type,
			Quota:  q.Quota,
			Usage:  uVal,
		})
	}

	return &projectUsage, nil
}

func (u *ProjectUsecase) UpdateProject(request *dtos.UpdateProjectRequest, userID uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !u.isProjAdminOrOrgAdmin(project, userID) {
		return nil, apiError.NewForbiddenError("user is not a project admin or organization admin")
	}

	// Update only provided fields
	if request.Name != "" {
		project.Name = request.Name
	}
	if request.Description != "" {
		project.Description = request.Description
	}

	if err := u.projRepo.UpdateProject(project); err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	// Fetch updated project with associations
	updatedProject, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return updatedProject, nil
}

func (u *ProjectUsecase) DeleteProject(projectID uuid.UUID, userID uuid.UUID) error {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	if !u.isProjAdminOrOrgAdmin(project, userID) {
		return apiError.NewForbiddenError("user is not a project admin or organization admin")
	}

	// Check if project has any active namespaces
	namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(projectID)
	if err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	if len(namespaces) > 0 {
		return apiError.NewBadRequestError("cannot delete project with existing namespaces")
	}

	// Check if project has any project quotas
	hasProjectQuotas, err := u.projRepo.HasProjectQuotas(projectID)
	if err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	if hasProjectQuotas {
		return apiError.NewBadRequestError("cannot delete project with existing project quotas. Please delete all project quotas first")
	}

	if err := u.projRepo.DeleteProject(projectID); err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	return nil
}

func (u *ProjectUsecase) AddAdmins(request *dtos.AddAdminsRequest, userID uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	org, err := u.orgRepo.GetOrganizationByID(project.OrganizationID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !u.isProjAdminOrOrgAdmin(project, userID) {
		return nil, apiError.NewForbiddenError("user is not a project admin or organization admin")
	}

	existing := make(map[uuid.UUID]struct{})
	for _, a := range project.Admins {
		existing[a.ID] = struct{}{}
	}

	seenReq := make(map[uuid.UUID]struct{})
	for _, adminID := range request.Admins {
		// Must be organization member to become project admin
		if !helper.ContainsUserID(org.Members, adminID) && !helper.ContainsUserID(org.Admins, adminID) {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("user %s is not a member of the organization", adminID))
		}
		if _, found := existing[adminID]; found {
			return nil, apiError.NewConflictError(fmt.Sprintf("user %s is already a project admin", adminID))
		}
		if _, found := seenReq[adminID]; found {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", adminID))
		}
		seenReq[adminID] = struct{}{}
		if _, err := u.userRepo.GetByID(adminID); err != nil {
			return nil, apiError.NewNotFoundError(fmt.Sprintf("user %s not found", adminID))
		}
	}

	users, err := u.userRepo.GetByIDs(request.Admins)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	project.Admins = append(project.Admins, users...)

	if err := u.projRepo.UpdateAdmins(project); err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return project, nil
}

func (u *ProjectUsecase) RemoveAdmins(request *dtos.RemoveAdminsRequest, userID uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !u.isProjAdminOrOrgAdmin(project, userID) {
		return nil, apiError.NewForbiddenError("user is not a project admin or organization admin")
	}

	// Create a map of admins to remove for quick lookup
	removeMap := make(map[uuid.UUID]struct{})
	seenReq := make(map[uuid.UUID]struct{})
	for _, adminID := range request.Admins {
		if _, found := seenReq[adminID]; found {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", adminID))
		}
		seenReq[adminID] = struct{}{}
		removeMap[adminID] = struct{}{}
	}

	// Check if admins to remove exist in project
	existing := make(map[uuid.UUID]struct{})
	for _, a := range project.Admins {
		existing[a.ID] = struct{}{}
	}

	for _, adminID := range request.Admins {
		if _, found := existing[adminID]; !found {
			return nil, apiError.NewNotFoundError(fmt.Sprintf("user %s is not an admin of this project", adminID))
		}
	}

	// Filter out admins to remove
	newAdmins := []models.User{}
	for _, admin := range project.Admins {
		if _, shouldRemove := removeMap[admin.ID]; !shouldRemove {
			newAdmins = append(newAdmins, admin)
		}
	}

	// Prevent removing all admins
	if len(newAdmins) == 0 {
		return nil, apiError.NewBadRequestError("cannot remove all admins from project. At least one admin must remain")
	}

	project.Admins = newAdmins

	if err := u.projRepo.UpdateAdmins(project); err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return project, nil
}

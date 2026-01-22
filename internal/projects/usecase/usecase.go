package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/ClearingHouse/internal/projects/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type ProjectUsecase struct {
	projRepo interfaces.ProjectRepository
	orgRepo  orgInterfaces.OrganizationRepository
	userRepo userInterfaces.UsersRepository
}

func NewProjectUsecase(projRepo interfaces.ProjectRepository, orgRepo orgInterfaces.OrganizationRepository, userRepo userInterfaces.UsersRepository) interfaces.ProjectUsecase {
	return &ProjectUsecase{
		projRepo: projRepo,
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

func (u *ProjectUsecase) CreateProject(request *dtos.CreateProjectRequest, userID uuid.UUID) error {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(org.Admins, userID) {
		return apiError.NewUnauthorizedError("user is not project admin")
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

func (u *ProjectUsecase) GetAllProjects() ([]models.Project, error) {
	projects, err := u.projRepo.GetAllProjects()
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}
	return projects, nil
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

	if !helper.ContainsUserID(project.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project admin")
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

	if !helper.ContainsUserID(project.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project admin")
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
		// Prevent removing admins through member removal
		if helper.ContainsUserID(project.Admins, memberID) {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("user %s is an admin and cannot be removed as a member directly", memberID))
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

	return project, nil
}

func (u *ProjectUsecase) GetAllUserProjects(userID uuid.UUID) ([]models.Project, error) {
	projects, err := u.projRepo.GetAllProjectsByUserID(userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return projects, nil
}

func (u *ProjectUsecase) GetProjectsByOrganizationID(orgID uuid.UUID, userID uuid.UUID) ([]models.Project, error) {
	org, err := u.orgRepo.GetOrganizationByID(orgID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(org.Members, userID) && !helper.ContainsUserID(org.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not a member of this organization")
	}

	projects, err := u.projRepo.GetProjectsByOrganizationID(orgID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}
	return projects, nil
}

func (u *ProjectUsecase) GetProjectMembers(projectID uuid.UUID, userID uuid.UUID) ([]models.User, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(project.Members, userID) && !helper.ContainsUserID(project.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project member")
	}

	members, err := u.projRepo.GetProjectMembers(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return members, nil
}

func (u *ProjectUsecase) GetProjectByID(projectID uuid.UUID, userID uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(project.Members, userID) && !helper.ContainsUserID(project.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project member")
	}

	return project, nil
}

func (u *ProjectUsecase) GetProjectUsage(projectID uuid.UUID, userID uuid.UUID) (*dtos.ProjectUsageResponse, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(project.Members, userID) && !helper.ContainsUserID(project.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project member")
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

	// Only admins can update project
	if !helper.ContainsUserID(project.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project admin")
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

	// Only admins can delete project
	if !helper.ContainsUserID(project.Admins, userID) {
		return apiError.NewUnauthorizedError("user is not project admin")
	}

	if err := u.projRepo.DeleteProject(projectID); err != nil {
		return apiError.NewInternalServerError(err.Error())
	}

	return nil
}

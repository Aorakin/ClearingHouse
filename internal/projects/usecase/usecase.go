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

// SERIOUS
func (u *ProjectUsecase) GetAllUserProjects(userID uuid.UUID) ([]models.Project, error) {
	projects, err := u.projRepo.GetAllProjectsByUserID(userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return projects, nil
}

func (u *ProjectUsecase) GetProjectByID(projectID uuid.UUID, userID uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	if !helper.ContainsUserID(project.Members, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project member")
	}

	return project, nil
}

func (u *ProjectUsecase) GetProjectQuota(projectID uuid.UUID) (*dtos.ProjectQuotaResponse, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return &dtos.ProjectQuotaResponse{
		ProjectID: project.ID,
		Quota:     1,
	}, nil
}

func (u *ProjectUsecase) GetProjectUsage(projectID uuid.UUID) (*dtos.ProjectUsageResponse, error) {
	project, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err.Error())
	}

	return &dtos.ProjectUsageResponse{
		ProjectID: project.ID,
		Usage:     1,
	}, nil
}

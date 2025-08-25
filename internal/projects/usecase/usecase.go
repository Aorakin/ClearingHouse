package usecase

import (
	"errors"
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/ClearingHouse/internal/projects/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
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
func (u *ProjectUsecase) CreateProject(request *dtos.CreateProjectRequest) error {
	org, err := u.orgRepo.GetOrganizationByID(request.OrganizationID)
	if err != nil {
		return err
	}

	if !helper.ContainsUserID(org.Admins, request.CreatorID) {
		return errors.New("unauthorized")
	}

	admin, err := u.userRepo.GetByID(request.CreatorID)
	if err != nil {
		return err
	}

	project := &models.Project{
		Name:           request.Name,
		Description:    request.Description,
		OrganizationID: request.OrganizationID,
		Admins:         []models.User{*admin},
	}

	return u.projRepo.CreateProject(project)
}
func (u *ProjectUsecase) GetAllProjects() ([]models.Project, error) {
	projects, err := u.projRepo.FindAllProjects()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (u *ProjectUsecase) AddMembers(projectID uuid.UUID, creator uuid.UUID, memberIDs []uuid.UUID) (*models.Project, error) {
	project, err := u.projRepo.FindProjectByID(projectID)
	if err != nil {
		return nil, err
	}

	org, err := u.orgRepo.GetOrganizationByID(project.OrganizationID)
	if err != nil {
		return nil, err
	}

	if !helper.ContainsUserID(project.Admins, creator) {
		return nil, fmt.Errorf("only project admins can add members")
	}

	existing := make(map[uuid.UUID]struct{})
	for _, m := range project.Members {
		existing[m.ID] = struct{}{}
	}

	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range memberIDs {
		// must be org member
		if !helper.ContainsUserID(org.Members, memberID) {
			return nil, fmt.Errorf("user %s is not a member of the organization", memberID)
		}
		if _, found := existing[memberID]; found {
			return nil, fmt.Errorf("user %s is already a project member", memberID)
		}
		if _, found := seenReq[memberID]; found {
			return nil, fmt.Errorf("duplicate user %s in request", memberID)
		}
		seenReq[memberID] = struct{}{}
		if _, err := u.userRepo.GetByID(memberID); err != nil {
			return nil, fmt.Errorf("user %s not found", memberID)
		}
	}

	users, err := u.userRepo.GetByIDs(memberIDs)
	if err != nil {
		return nil, err
	}

	project.Members = append(project.Members, users...)

	if err := u.projRepo.UpdateMembers(project); err != nil {
		return nil, err
	}

	return project, nil
}

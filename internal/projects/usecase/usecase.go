package usecase

import (
	"errors"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dto"
	"github.com/ClearingHouse/internal/projects/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	"github.com/google/uuid"
)

type ProjectsUsecase struct {
	projectsRepository interfaces.ProjectsRepository
	usersRepository    userInterfaces.UsersRepository
}

func NewProjectsUsecase(projectsRepository interfaces.ProjectsRepository, usersRepository userInterfaces.UsersRepository) interfaces.ProjectsUsecase {
	return &ProjectsUsecase{
		projectsRepository: projectsRepository,
		usersRepository:    usersRepository,
	}
}

func (u *ProjectsUsecase) GetAll(userID uuid.UUID) ([]models.Project, error) {
	projects, err := u.projectsRepository.GetAll(userID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (u *ProjectsUsecase) GetByID(userID uuid.UUID, projectID uuid.UUID) (*models.Project, error) {
	project, err := u.projectsRepository.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	for _, owner := range project.Owners {
		if owner.ID == userID {
			return project, nil
		}
	}
	for _, member := range project.Members {
		if member.ID == userID {
			return project, nil
		}
	}

	return nil, errors.New("unauthorized: user is not an owner or member of the project")
}

func (u *ProjectsUsecase) Create(userID uuid.UUID, projectRequest *dto.CreateProjectRequest) (*models.Project, error) {
	var owners []*models.User
	user, err := u.usersRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	owners = append(owners, user)

	project := models.Project{
		Title:       projectRequest.Title,
		Description: projectRequest.Description,
		Owners:      owners,
	}

	createdProject, err := u.projectsRepository.Create(&project)
	if err != nil {
		return nil, err
	}

	return createdProject, nil
}

func (u *ProjectsUsecase) Update(userID uuid.UUID, projectID uuid.UUID, projectRequest *dto.CreateProjectRequest) (*models.Project, error) {
	project, err := u.projectsRepository.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	isAuthorized := false
	for _, owner := range project.Owners {
		if owner.ID == userID {
			isAuthorized = true
			break
		}
	}
	if !isAuthorized {
		for _, member := range project.Members {
			if member.ID == userID {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return nil, errors.New("unauthorized: user is not an owner or member of the project")
	}

	project.Title = projectRequest.Title
	project.Description = projectRequest.Description

	updatedProject, err := u.projectsRepository.Update(projectID, project)
	if err != nil {
		return nil, err
	}

	return updatedProject, nil
}

func (u *ProjectsUsecase) Delete(userID uuid.UUID, projectID uuid.UUID) error {
	project, err := u.projectsRepository.GetByID(projectID)
	if err != nil {
		return err
	}

	isOwner := false
	for _, owner := range project.Owners {
		if owner.ID == userID {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return errors.New("unauthorized: user is not an owner of the project")
	}

	err = u.projectsRepository.Delete(projectID)
	if err != nil {
		return err
	}

	return nil
}

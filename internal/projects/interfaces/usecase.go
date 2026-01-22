package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	GetAllProjects() ([]models.Project, error)
	CreateProject(request *dtos.CreateProjectRequest, userID uuid.UUID) error
	AddMembers(request *dtos.AddMembersRequest, userID uuid.UUID) (*models.Project, error)
	RemoveMembers(request *dtos.RemoveMembersRequest, userID uuid.UUID) (*models.Project, error)

	GetAllUserProjects(userID uuid.UUID) ([]models.Project, error)
	GetProjectByID(projectID uuid.UUID, userID uuid.UUID) (*models.Project, error)
	GetProjectsByOrganizationID(orgID uuid.UUID, userID uuid.UUID) ([]models.Project, error)
	GetProjectMembers(projectID uuid.UUID, userID uuid.UUID) ([]models.User, error)

	GetProjectUsage(projectID uuid.UUID, userID uuid.UUID) (*dtos.ProjectUsageResponse, error)
	UpdateProject(request *dtos.UpdateProjectRequest, userID uuid.UUID) (*models.Project, error)
	DeleteProject(projectID uuid.UUID, userID uuid.UUID) error
}

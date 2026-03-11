package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	projectDtos "github.com/ClearingHouse/internal/projects/dtos"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	GetAllProjects(userID uuid.UUID, isSuperAdmin bool) ([]projectDtos.ProjectResponse, error)
	CreateProject(request *projectDtos.CreateProjectRequest, userID uuid.UUID) error
	AddMembers(request *projectDtos.AddMembersRequest, userID uuid.UUID) (*models.Project, error)
	RemoveMembers(request *projectDtos.RemoveMembersRequest, userID uuid.UUID) (*models.Project, error)
	AddAdmins(request *projectDtos.AddAdminsRequest, userID uuid.UUID) (*models.Project, error)
	RemoveAdmins(request *projectDtos.RemoveAdminsRequest, userID uuid.UUID) (*models.Project, error)

	GetAllUserProjects(userID uuid.UUID) ([]projectDtos.ProjectResponse, error)
	GetProjectByID(projectID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) (*projectDtos.ProjectResponse, error)
	GetProjectsByOrganizationID(orgID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) ([]projectDtos.ProjectResponse, error)
	GetProjectMembers(projectID uuid.UUID, userID uuid.UUID, isSuperAdmin bool) ([]models.User, error)

	GetProjectUsage(projectID uuid.UUID, userID uuid.UUID) (*projectDtos.ProjectUsageResponse, error)
	UpdateProject(request *projectDtos.UpdateProjectRequest, userID uuid.UUID) (*models.Project, error)
	DeleteProject(projectID uuid.UUID, userID uuid.UUID) error
}

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

	GetAllUserProjects(userID uuid.UUID) ([]models.Project, error)
	GetProjectByID(projectID uuid.UUID, userID uuid.UUID) (*models.Project, error)

	GetProjectUsage(projectID uuid.UUID, userID uuid.UUID) (*dtos.ProjectUsageResponse, error)
}

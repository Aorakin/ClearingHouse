package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	CreateProject(request *dtos.CreateProjectRequest) error
	GetAllProjects() ([]models.Project, error)
	AddMembers(request *dtos.AddMembersRequest) (*models.Project, error)

	GetAllUserProjects(userID uuid.UUID) ([]models.Project, error)
	GetProjectByID(projectID uuid.UUID, userID uuid.UUID) (*models.Project, apiError.ApiErr)
}

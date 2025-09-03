package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetAllProjects() ([]models.Project, error)

	GetProjectByID(id uuid.UUID) (*models.Project, error)
	UpdateMembers(project *models.Project) error

	GetAllProjectsByUserID(userID uuid.UUID) ([]models.Project, error)
	// GetProjectQuota(projectID uuid.UUID) (*models.Project, error)
	// GetProjectUsage(projectID uuid.UUID) (*models.Project, error)
}

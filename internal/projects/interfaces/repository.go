package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	FindAllProjects() ([]models.Project, error)
	FindProjectByID(id uuid.UUID) (*models.Project, error)
	UpdateMembers(project *models.Project) error
}

package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type ProjectsRepository interface {
	GetAll(userID uuid.UUID) ([]models.Project, error)
	GetByID(id uuid.UUID) (*models.Project, error)
	Create(project *models.Project) (*models.Project, error)
	Update(id uuid.UUID, project *models.Project) (*models.Project, error)
	Delete(id uuid.UUID) error
}

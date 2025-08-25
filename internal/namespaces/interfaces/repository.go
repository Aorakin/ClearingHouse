package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type NamespaceRepository interface {
	Create(namespace *models.Namespace) error
	GetAll() ([]models.Namespace, error)
	GetByID(id uuid.UUID) (*models.Namespace, error)
	UpdateMembers(namespace *models.Namespace) error
	FindAllNamespacesByProjectID(projectID uuid.UUID) ([]models.Namespace, error)
}

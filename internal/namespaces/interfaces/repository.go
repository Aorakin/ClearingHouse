package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type NamespacesRepository interface {
	GetById(ID uuid.UUID) (*models.Namespace, error)
	Create(namespace *models.Namespace) (*models.Namespace, error)
	Update(namespace *models.Namespace) (*models.Namespace, error)
}

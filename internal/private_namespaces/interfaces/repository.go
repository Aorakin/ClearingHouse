package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type PrivateNamespaceRepository interface {
	GetPrivateNamespaceByOwnerID(ownerID uuid.UUID) (*models.Namespace, error)
	CreatePrivateNamespace(namespace *models.Namespace) error
}

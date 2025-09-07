package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type NamespaceRepository interface {
	Create(namespace *models.Namespace) error
	GetAll() ([]models.Namespace, error)
	GetNamespaceByID(namespaceID uuid.UUID) (*models.Namespace, error)
	UpdateNamespace(namespace *models.Namespace) error
	UpdateMembers(namespace *models.Namespace) error

	GetAllNamespacesByProjectID(projectID uuid.UUID) ([]models.Namespace, error)
	GetAllNamespacesByUserID(userID uuid.UUID) ([]models.Namespace, error)

	GetNamespaceQuotas(namespaceID uuid.UUID) ([]models.NamespaceQuota, error)
	GetNamespaceTickets(namespaceID, resourcePoolID, quotaID uuid.UUID) ([]models.Ticket, error)
}

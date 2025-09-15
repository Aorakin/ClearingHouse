package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
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
	GetAllNamespacesByProjectAndUserID(projectID, userID uuid.UUID) ([]models.Namespace, error)

	GetNamespaceQuotas(namespaceID uuid.UUID) ([]models.NamespaceQuota, error)
	GetNamespaceTickets(namespaceID, resourcePoolID, quotaID uuid.UUID) ([]models.Ticket, error)

	GetNamespaceQuotaByType(namespaceID uuid.UUID) (*dtos.ResourceQuotaResponse, error)
	GetNamespaceUsageByType(namespaceID uuid.UUID) (*dtos.ResourceUsageResponse, error)
}

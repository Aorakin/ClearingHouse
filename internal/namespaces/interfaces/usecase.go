package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

type NamespacesUsecase interface {
	Create(projectID uuid.UUID) (*models.Namespace, error)
	AddMember(namespaceID uuid.UUID, memberIDs uuid.UUIDs) (*models.Namespace, error)
	ChangeQuota(namespaceID uuid.UUID, quotaID uuid.UUID) (*models.Namespace, error)
}

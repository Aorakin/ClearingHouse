package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/google/uuid"
)

type NamespaceUsecase interface {
	CreateNamespace(request *dtos.CreateNamespaceRequest, userID uuid.UUID) (*models.Namespace, error)
	GetAllNamespaces() ([]models.Namespace, error)
	AddMembers(request *dtos.AddMembersRequest, userID uuid.UUID) (*models.Namespace, error)

	GetAllUserNamespaces(userID uuid.UUID) ([]models.Namespace, error)
	GetNamespace(namespaceID uuid.UUID, userID uuid.UUID) (*models.Namespace, error)
}

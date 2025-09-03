package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type NamespaceUsecase interface {
	CreateNamespace(request *dtos.CreateNamespaceRequest) (*models.Namespace, error)
	GetAllNamespaces() ([]models.Namespace, error)
	AddMembers(request *dtos.AddMembersRequest) (*models.Namespace, error)

	GetAllUserNamespaces(userID uuid.UUID) ([]models.Namespace, apiError.ApiErr)
	GetNamespace(namespaceID uuid.UUID, userID uuid.UUID) (*models.Namespace, apiError.ApiErr)
}

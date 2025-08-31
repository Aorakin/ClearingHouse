package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
)

type NamespaceUsecase interface {
	CreateNamespace(request *dtos.CreateNamespaceRequest) (*models.Namespace, error)
	GetAllNamespaces() ([]models.Namespace, error)
	AddMembers(request *dtos.AddMembersRequest) (*models.Namespace, error)
}

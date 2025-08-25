package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
)

type NamespaceUsecase interface {
	CreateNamespace(namespace *dtos.CreateNamespaceRequest) error
	GetAllNamespaces() ([]models.Namespace, error)
}

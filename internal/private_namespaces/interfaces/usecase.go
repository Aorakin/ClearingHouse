package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/private_namespaces/dtos"
	"github.com/google/uuid"
)

type PrivateNamespaceUsecase interface {
	GetPrivateNamespaceByOwnerID(ownerID uuid.UUID) (*models.Namespace, error)
	CreatePrivateNamespace(request *dtos.CreatePrivateNamespaceRequest, userID uuid.UUID) (*models.Namespace, error)
	CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error)
}

package usecase

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/private_namespaces/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivateNamespaceUsecase struct {
	privNamespaceRepository interfaces.PrivateNamespaceRepository
}

func NewPrivateNamespaceUsecase(privNamespaceRepository interfaces.PrivateNamespaceRepository) interfaces.PrivateNamespaceUsecase {
	return &PrivateNamespaceUsecase{
		privNamespaceRepository: privNamespaceRepository,
	}
}

func (u *PrivateNamespaceUsecase) GetPrivateNamespaceByOwnerID(ownerID uuid.UUID) (*models.Namespace, error) {
	namespace, err := u.privNamespaceRepository.GetPrivateNamespaceByOwnerID(ownerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apiError.NewNotFoundError("namespace not found")
		}
		return nil, err
	}
	return namespace, nil
}

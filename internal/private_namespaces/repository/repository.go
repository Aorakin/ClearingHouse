package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/private_namespaces/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivateNamespaceRepository struct {
	db *gorm.DB
}

func NewPrivateNamespaceRepository(db *gorm.DB) interfaces.PrivateNamespaceRepository {
	return &PrivateNamespaceRepository{
		db: db,
	}
}

func (r *PrivateNamespaceRepository) GetPrivateNamespaceByOwnerID(ownerID uuid.UUID) (*models.Namespace, error) {
	var user models.User
	err := r.db.Preload("Namespace").First(&user, "id = ?", ownerID).Error
	if err != nil {
		return nil, err
	}
	return user.Namespace, nil
}

func (r *PrivateNamespaceRepository) CreatePrivateNamespace(namespace *models.Namespace) error {
	return r.db.Create(namespace).Error
}

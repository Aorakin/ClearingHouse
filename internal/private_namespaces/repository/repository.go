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
	var namespace models.Namespace
	err := r.db.First(&namespace, "owner_id = ?", ownerID).Error
	if err != nil {
		return nil, err
	}
	return &namespace, nil
}

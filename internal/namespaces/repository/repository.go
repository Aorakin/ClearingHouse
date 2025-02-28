package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NamespacesRepository struct {
	db *gorm.DB
}

func NewNamespacesRepository(db *gorm.DB) interfaces.NamespacesRepository {
	return &NamespacesRepository{db: db}
}

func (r *NamespacesRepository) GetById(ID uuid.UUID) (*models.Namespace, error) {
	var namespace models.Namespace
	err := r.db.First(&namespace, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}

	return &namespace, nil
}

func (r *NamespacesRepository) Create(namespace *models.Namespace) (*models.Namespace, error) {
	err := r.db.Create(namespace).Error
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (r *NamespacesRepository) Update(namespace *models.Namespace) (*models.Namespace, error) {
	err := r.db.Model(&models.Namespace{}).Where("id = ?", namespace.ID).Updates(namespace).Error
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

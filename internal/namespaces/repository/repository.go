package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NamespaceRepository struct {
	db *gorm.DB
}

func NewNamespaceRepository(db *gorm.DB) interfaces.NamespaceRepository {
	return &NamespaceRepository{
		db: db,
	}
}

func (r *NamespaceRepository) Create(namespace *models.Namespace) error {
	return r.db.Create(namespace).Error
}

func (r *NamespaceRepository) GetAll() ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) GetByID(id uuid.UUID) (*models.Namespace, error) {
	var namespace models.Namespace
	err := r.db.Preload("Members").First(&namespace, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &namespace, nil
}

func (r *NamespaceRepository) FindAllNamespacesByProjectID(projectID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Where("project_id = ?", projectID).Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) UpdateMembers(namespace *models.Namespace) error {
	return r.db.Model(namespace).Association("Members").Replace(namespace.Members)
}

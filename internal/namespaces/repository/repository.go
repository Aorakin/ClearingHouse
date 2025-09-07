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

func (r *NamespaceRepository) GetNamespaceByID(namespaceID uuid.UUID) (*models.Namespace, error) {
	var namespace models.Namespace
	err := r.db.Preload("Members").First(&namespace, "id = ?", namespaceID).Error
	if err != nil {
		return nil, err
	}
	return &namespace, nil
}

func (r *NamespaceRepository) GetAllNamespacesByProjectID(projectID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Where("project_id = ?", projectID).Find(&namespaces).Error
	return namespaces, err
}

func (r *NamespaceRepository) UpdateMembers(namespace *models.Namespace) error {
	return r.db.Model(namespace).Association("Members").Replace(namespace.Members)
}

func (r *NamespaceRepository) UpdateNamespace(namespace *models.Namespace) error {
	return r.db.Save(namespace).Error
}

func (r *NamespaceRepository) GetAllNamespacesByUserID(userID uuid.UUID) ([]models.Namespace, error) {
	var user models.User
	if err := r.db.Debug().Preload("MemberNamespaces").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.MemberNamespaces, nil
}

func (r *NamespaceRepository) GetNamespaceQuotas(namespaceID uuid.UUID) ([]models.NamespaceQuota, error) {
	var quotas []models.NamespaceQuota
	err := r.db.Preload("Resources.ResourceProperties").Joins("JOIN namespace_quotas nq ON nq.namespace_id = ?", namespaceID).Find(&quotas).Error
	if err != nil {
		return nil, err
	}
	return quotas, nil
}

func (r *NamespaceRepository) GetNamespaceTickets(namespaceID, resourcePoolID, quotaID uuid.UUID) ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.db.Debug().
		Preload("Resources").
		Where("namespace_id = ? AND resource_pool_id = ? AND quota_id = ?", namespaceID, resourcePoolID, quotaID).
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *NamespaceRepository) GetAllNamespacesByProjectAndUserID(projectID, userID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	err := r.db.Joins("JOIN namespace_members nm ON nm.namespace_id = namespaces.id").
		Where("namespaces.project_id = ? AND nm.user_id = ?", projectID, userID).
		Find(&namespaces).Error
	return namespaces, err
}

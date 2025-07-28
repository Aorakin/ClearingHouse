package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) interfaces.OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) CreateOrganization(org *models.Organization) (*models.Organization, error) {
	if err := r.db.Create(org).Error; err != nil {
		return nil, err
	}
	return org, nil
}
func (r *OrganizationRepository) GetOrganizationByID(id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}
func (r *OrganizationRepository) UpdateOrganization(org *models.Organization) (*models.Organization, error) {
	if err := r.db.Save(org).Error; err != nil {
		return nil, err
	}
	return org, nil
}
func (r *OrganizationRepository) DeleteOrganization(id uuid.UUID) error {
	if err := r.db.Delete(&models.Organization{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *OrganizationRepository) GetOrganizations() ([]models.Organization, error) {
	var organizations []models.Organization
	if err := r.db.Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

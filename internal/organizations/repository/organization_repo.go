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
	if err := r.db.
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Projects", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("ResourcePools", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("Quotas", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("GivenQuotas", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) GetOrganizationByDomain(domain string) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Where("domain = ? AND domain != ''", domain).First(&org).Error; err != nil {
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		org := models.Organization{BaseModel: models.BaseModel{ID: id}}

		if err := tx.Model(&org).Association("Members").Clear(); err != nil {
			return err
		}

		if err := tx.Model(&org).Association("Admins").Clear(); err != nil {
			return err
		}

		return tx.Delete(&org).Error
	})
}

func (r *OrganizationRepository) GetOrganizations() ([]models.Organization, error) {
	var organizations []models.Organization
	if err := r.db.
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Projects", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("ResourcePools", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("Quotas", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("GivenQuotas", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Order("name").
		Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *OrganizationRepository) GetOrganizationsByAdminOrProjectAdmin(userID uuid.UUID) ([]models.Organization, error) {
	var organizations []models.Organization

	subqueryOrgAdmin := r.db.Table("organization_admins").
		Select("organization_id").
		Where("user_id = ?", userID)

	subqueryProjAdmin := r.db.Table("project_admins").
		Joins("JOIN projects ON projects.id = project_admins.project_id").
		Select("DISTINCT projects.organization_id").
		Where("project_admins.user_id = ? AND projects.deleted_at IS NULL", userID)

	if err := r.db.
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Projects", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("ResourcePools", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("Quotas", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Preload("GivenQuotas", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).
		Where("id IN (?) OR id IN (?)", subqueryOrgAdmin, subqueryProjAdmin).
		Order("name").
		Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *OrganizationRepository) IsProjectAdminInOrg(orgID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Table("project_admins").
		Joins("JOIN projects ON projects.id = project_admins.project_id").
		Where("projects.organization_id = ? AND project_admins.user_id = ? AND projects.deleted_at IS NULL", orgID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrganizationRepository) UpdateMembers(org *models.Organization) error {
	return r.db.Model(org).Association("Members").Replace(org.Members)
}

func (r *OrganizationRepository) UpdateAdmins(org *models.Organization) error {
	return r.db.Model(org).Association("Admins").Replace(org.Admins)
}

func (r *OrganizationRepository) GetMembers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Order("email").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *OrganizationRepository) GetOrganizationMembers(orgID uuid.UUID) ([]models.User, error) {
	var org models.Organization
	if err := r.db.Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).First(&org, "id = ?", orgID).Error; err != nil {
		return nil, err
	}
	return org.Members, nil
}

package usecase

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/google/uuid"
)

type OrganizationUsecase struct {
	orgRepo interfaces.OrganizationRepository
}

func NewOrganizationUsecase(orgRepo interfaces.OrganizationRepository) interfaces.OrganizationUsecase {
	return &OrganizationUsecase{
		orgRepo: orgRepo,
	}
}

func (u *OrganizationUsecase) CreateOrganization(org *dtos.CreateOrganization) (*models.Organization, error) {
	organization := &models.Organization{
		Name:        org.Name,
		Description: org.Description,
	}
	return u.orgRepo.CreateOrganization(organization)
}

func (u *OrganizationUsecase) GetOrganizationByID(id uuid.UUID) (*models.Organization, error) {
	return u.orgRepo.GetOrganizationByID(id)
}

func (u *OrganizationUsecase) UpdateOrganization(id uuid.UUID, org *dtos.UpdateOrganization) (*models.Organization, error) {
	organization := &models.Organization{
		BaseModel:   models.BaseModel{ID: id},
		Name:        org.Name,
		Description: org.Description,
	}
	return u.orgRepo.UpdateOrganization(organization)
}

func (u *OrganizationUsecase) DeleteOrganization(id uuid.UUID) error {
	return u.orgRepo.DeleteOrganization(id)
}

func (u *OrganizationUsecase) GetOrganizations() ([]models.Organization, error) {
	return u.orgRepo.GetOrganizations()
}

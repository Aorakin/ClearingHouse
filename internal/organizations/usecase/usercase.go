package usecase

import (
	"errors"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	"github.com/google/uuid"
)

type OrganizationUsecase struct {
	orgRepo  interfaces.OrganizationRepository
	userRepo userInterfaces.UsersRepository
}

func NewOrganizationUsecase(orgRepo interfaces.OrganizationRepository, userRepo userInterfaces.UsersRepository) interfaces.OrganizationUsecase {
	return &OrganizationUsecase{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

func (u *OrganizationUsecase) CreateOrganization(org *dtos.CreateOrganization) (*models.Organization, error) {
	creator, err := u.userRepo.GetByID(org.Creator)
	if err != nil {
		return nil, err
	}

	organization := &models.Organization{
		Name:        org.Name,
		Description: org.Description,
		Admins:      []models.User{*creator},
		Members:     []models.User{*creator},
	}
	return u.orgRepo.CreateOrganization(organization)
}

func (u *OrganizationUsecase) GetOrganizationByID(id uuid.UUID, userID uuid.UUID) (*models.Organization, error) {
	organization, err := u.orgRepo.GetOrganizationByID(id)
	if err != nil {
		return nil, err
	}

	if !helper.ContainsUserID(organization.Members, userID) {
		return nil, errors.New("unauthorized")
	}

	return organization, nil
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

package usecase

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dto"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"github.com/google/uuid"

	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
)

type NamespacesUsecase struct {
	namespacesRepository interfaces.NamespacesRepository
	usersRepository      userInterfaces.UsersRepository
}

func NewNamespacesUsecase(namespacesRepository interfaces.NamespacesRepository, usersRepository userInterfaces.UsersRepository) interfaces.NamespacesUsecase {
	return &NamespacesUsecase{
		namespacesRepository: namespacesRepository,
		usersRepository:      usersRepository,
	}
}

// ADD AUTHENTICATION AND AUTHORIZATION

func (u *NamespacesUsecase) GetAll(projectID uuid.UUID) ([]models.Namespace, error) {
	var namespaces []models.Namespace
	return namespaces, nil
}

func (u *NamespacesUsecase) GetByID(namespaceID uuid.UUID) (*models.Namespace, error) {
	var namespace *models.Namespace
	return namespace, nil
}

func (u *NamespacesUsecase) Create(projectID uuid.UUID, request *dto.CreateNamespaceRequest) (*models.Namespace, error) {
	namespace := models.Namespace{
		ProjectID: projectID,
	}

	createdNamespace, err := u.namespacesRepository.Create(&namespace)
	if err != nil {
		return nil, err
	}

	return createdNamespace, nil
}

func (u *NamespacesUsecase) AddMembers(userID uuid.UUID, namespaceID uuid.UUID, memberIDs uuid.UUIDs) (*models.Namespace, error) {
	namespace, err := u.namespacesRepository.GetById(namespaceID)
	if err != nil {
		return nil, err
	}

	var newMembers []*models.User
	for _, id := range memberIDs {
		user, err := u.usersRepository.GetUser(id)
		if err != nil {
			return nil, err
		}
		newMembers = append(newMembers, user)
	}

	namespace.Members = append(namespace.Members, newMembers...)

	updatedNamespace, err := u.namespacesRepository.Update(namespace)
	if err != nil {
		return nil, err
	}

	return updatedNamespace, nil
}

func (u *NamespacesUsecase) ChangeQuota(userID uuid.UUID, namespaceID uuid.UUID, quotaID uuid.UUID) (*models.Namespace, error) {
	namespace, err := u.namespacesRepository.GetById(namespaceID)
	if err != nil {
		return nil, err
	}

	// ADD CHECK QUOTA EXIST

	namespace.QuotaID = quotaID

	updatedNamespace, err := u.namespacesRepository.Update(namespace)
	if err != nil {
		return nil, err
	}

	return updatedNamespace, nil
}

func (u *NamespacesUsecase) RequestTicket(userID uuid.UUID, namespaceID uuid.UUID, request *dto.CreateTicketRequest) (*models.Ticket, error) {
	var ticket models.Ticket
	return &ticket, nil
}

func (u *NamespacesUsecase) TerminateTicket(userID uuid.UUID, namespaceID uuid.UUID, ticketID uuid.UUID) error {
	return nil
}

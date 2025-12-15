package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	resourceInterfaces "github.com/ClearingHouse/internal/resources/interfaces"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type TicketUsecase struct {
	namespaceRepo namespaceInterfaces.NamespaceRepository
	ticketRepo    interfaces.TicketRepository
	quotaRepo     quotaInterfaces.QuotaRepository
	userRepo      userInterfaces.UsersRepository
	resourceRepo  resourceInterfaces.ResourceRepository
}

func NewTicketUsecase(namespaceRepo namespaceInterfaces.NamespaceRepository, ticketRepo interfaces.TicketRepository, quotaRepo quotaInterfaces.QuotaRepository, userRepo userInterfaces.UsersRepository, resourceRepo resourceInterfaces.ResourceRepository) interfaces.TicketUsecase {
	return &TicketUsecase{
		namespaceRepo: namespaceRepo,
		ticketRepo:    ticketRepo,
		quotaRepo:     quotaRepo,
		userRepo:      userRepo,
		resourceRepo:  resourceRepo,
	}
}

func (u *TicketUsecase) isNamespaceMember(namespaceID uuid.UUID, userID uuid.UUID) error {
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get namespace: %w", err))
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}

	if !helper.ContainsUserID(namespace.Members, user.ID) {
		return apiError.NewForbiddenError(fmt.Errorf("user is not a member of the namespace"))
	}

	return nil
}

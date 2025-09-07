package usecase

import (
	"fmt"
	"log"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type TicketUsecase struct {
	namespaceRepo namespaceInterfaces.NamespaceRepository
	ticketRepo    interfaces.TicketRepository
	quotaRepo     quotaInterfaces.QuotaRepository
}

func NewTicketUsecase(namespaceRepo namespaceInterfaces.NamespaceRepository, ticketRepo interfaces.TicketRepository, quotaRepo quotaInterfaces.QuotaRepository) interfaces.TicketUsecase {
	return &TicketUsecase{
		namespaceRepo: namespaceRepo,
		ticketRepo:    ticketRepo,
		quotaRepo:     quotaRepo,
	}
}

func (u *TicketUsecase) CreateTicket(request *dtos.CreateTicketRequest, userID uuid.UUID) (*models.Ticket, error) {
	if len(request.Resources) == 0 {
		return nil, apiError.NewBadRequestError("at least one resource is required")
	}

	isMember, err := u.isNamespaceMember(userID, request.NamespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !isMember {
		return nil, apiError.NewUnauthorizedError("user is not a member of the namespace")
	}

	isAssigned, err := u.quotaRepo.IsAssigned(request.NamespaceID, request.QuotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !isAssigned {
		return nil, apiError.NewUnauthorizedError("quota is not assigned to the namespace")
	}

	quota, err := u.quotaRepo.GetNamespaceQuotaByID(request.QuotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	quotaResources := make(map[uuid.UUID]models.ResourceQuantity)
	for _, resource := range quota.Resources {
		quotaResources[resource.ResourceProp.ResourceID] = resource
	}
	seenReq := make(map[uuid.UUID]struct{})
	totalCredit := float32(0)

	for _, resource := range request.Resources {
		if _, ok := quotaResources[resource.ResourceID]; !ok {
			return nil, apiError.NewForbiddenError("resource is not allowed")
		}
		resourceQuantity := quotaResources[resource.ResourceID]

		if _, ok := seenReq[resource.ResourceID]; ok {
			return nil, apiError.NewBadRequestError("resource is duplicated")
		}
		seenReq[resource.ResourceID] = struct{}{}

		usage, err := u.ticketRepo.GetResourceUsage(request.NamespaceID, request.QuotaID, resource.ResourceID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}

		log.Printf("Resource ID: %s usage %d", resource.ResourceID, usage)
		maxQuota := resourceQuantity.Quantity
		if usage+resource.Quantity > maxQuota {
			return nil, apiError.NewForbiddenError("namespace usage exceeds quota limit for resource")
		}

		if request.Duration <= 0 {
			return nil, apiError.NewBadRequestError("duration must be greater than 0")
		}
		if request.Duration > resourceQuantity.ResourceProp.MaxDuration {
			return nil, apiError.NewForbiddenError("duration exceeds max limit for resource")
		}

		totalCredit += float32(resource.Quantity) * resourceQuantity.ResourceProp.Price * request.Duration
	}

	namespace, err := u.namespaceRepo.GetNamespaceByID(request.NamespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	if totalCredit > namespace.Credit {
		return nil, apiError.NewForbiddenError(fmt.Errorf("not enough credit want to use %.2f but only have %.2f", totalCredit, namespace.Credit))
	}
	namespace.Credit -= totalCredit

	err = u.namespaceRepo.UpdateNamespace(namespace)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	ticket := &models.Ticket{
		NamespaceID:    request.NamespaceID,
		Name:           request.Name,
		OwnerID:        userID,
		QuotaID:        request.QuotaID,
		ResourcePoolID: quota.ResourcePoolID,
		Status:         "created",
	}
	err = u.ticketRepo.CreateTicket(ticket)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	for _, resource := range request.Resources {
		ticketResource := &models.TicketResource{
			ResourceID: resource.ResourceID,
			Quantity:   resource.Quantity,
			TicketID:   ticket.ID,
		}
		err = u.ticketRepo.CreateTicketResource(ticketResource)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}
	}

	// reload ticket with resources
	ticket, err = u.ticketRepo.GetTicketByID(ticket.ID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return ticket, nil
}

func (u *TicketUsecase) GetNamespaceTickets(namespaceID uuid.UUID, userID uuid.UUID) ([]models.Ticket, error) {
	isMember, err := u.isNamespaceMember(userID, namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !isMember {
		return nil, apiError.NewUnauthorizedError("user is not a member of the namespace")
	}

	tickets, err := u.ticketRepo.GetNamespaceTickets(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return tickets, nil
}

func (u *TicketUsecase) isNamespaceMember(userID uuid.UUID, namespaceID uuid.UUID) (bool, error) {
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return false, err
	}

	if !helper.ContainsUserID(namespace.Members, userID) {
		return false, nil
	}

	return true, nil
}

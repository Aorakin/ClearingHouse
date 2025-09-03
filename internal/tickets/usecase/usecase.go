package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/ClearingHouse/internal/tickets/interfaces"
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

func (u *TicketUsecase) CreateTicket(request *dtos.CreateTicketRequest) (*models.Ticket, error) {
	isMember, err := u.IsNamespaceMember(request.Creator, request.NamespaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify namespace membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of the namespace")
	}

	namespace, err := u.namespaceRepo.GetByID(request.NamespaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	if namespace.QuotaGroupID == nil {
		return nil, fmt.Errorf("namespace does not have a quota group assigned")
	}

	quota, err := u.quotaRepo.FindNamespaceQuotaGroupByID(*namespace.QuotaGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace quota group: %w", err)
	}

	// validate resources quantity
	allowedResources := make(map[uuid.UUID]struct{})
	for _, resource := range quota.Resources {
		allowedResources[resource.ResourceProperty.ResourceID] = struct{}{}
	}

	totalCredit := float32(0)
	seen := make(map[uuid.UUID]struct{})
	for _, resource := range request.Resources {
		// check if resource exist in quota
		if _, ok := allowedResources[resource.ResourceID]; !ok {
			return nil, fmt.Errorf("resource %s is not allowed", resource.ResourceID)
		}

		// check if resource is duplicated
		if _, ok := seen[resource.ResourceID]; ok {
			return nil, fmt.Errorf("resource %s is duplicated", resource.ResourceID)
		}
		seen[resource.ResourceID] = struct{}{}

		// check namespace usage
		namespaceUsage, err := u.ticketRepo.GetNamespaceUsage(request.NamespaceID, quota.ID, resource.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get namespace usage: %w", err)
		}

		maxQuota, err := u.quotaRepo.GetNamespaceQuotaQuantity(quota.ID, resource.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get max quota for resource %s: %w", resource.ResourceID, err)
		}

		if namespaceUsage+resource.Quantity > maxQuota {
			return nil, fmt.Errorf("namespace usage exceeds quota limit for resource %s", resource.ResourceID)
		}

		resourceProperty, err := u.quotaRepo.GetResourcePropertyByNamespace(*namespace.QuotaGroupID, resource.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource property for resource %s: %w", resource.ResourceID, err)
		}

		totalCredit += float32(resource.Quantity) * resourceProperty.Price * float32(request.Duration)
	}

	if totalCredit > namespace.Credit {
		return nil, fmt.Errorf("not enough credit")
	}

	namespace.Credit -= totalCredit
	err = u.namespaceRepo.UpdateNamespace(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to update namespace: %w", err)
	}

	ticket := &models.Ticket{
		NamespaceID: request.NamespaceID,
		Name:        request.Name,
		OwnerID:     request.Creator,
		Status:      "created",
	}

	err = u.ticketRepo.CreateTicket(ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	for _, resource := range request.Resources {
		ticketResource := &models.TicketResource{
			QuotaID:    *namespace.QuotaGroupID,
			ResourceID: resource.ResourceID,
			Quantity:   resource.Quantity,
			TicketID:   ticket.ID,
		}
		err = u.ticketRepo.CreateTicketResource(ticketResource)
		if err != nil {
			return nil, fmt.Errorf("failed to create ticket resource: %w", err)
		}
	}

	ticket, err = u.ticketRepo.GetTicketByID(ticket.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created ticket: %w", err)
	}

	return ticket, nil
}

func (u *TicketUsecase) IsNamespaceMember(userID uuid.UUID, namespaceID uuid.UUID) (bool, error) {
	namespace, err := u.namespaceRepo.GetByID(namespaceID)
	if err != nil {
		return false, err
	}

	if !helper.ContainsUserID(namespace.Members, userID) {
		return false, nil
	}

	return true, nil
}

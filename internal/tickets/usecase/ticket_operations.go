package usecase

import (
	"fmt"
	"time"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *TicketUsecase) validateCreateTicketRequest(request *dtos.CreateTicketRequest, quotaResourcesMap map[uuid.UUID]models.ResourceQuantity) (*models.Namespace, float32, error) {
	seenReq := make(map[uuid.UUID]struct{})
	totalCredit := float32(0)
	price := float32(0)

	if len(request.Resources) == 0 {
		return nil, 0, apiError.NewBadRequestError("at least one resource is required")
	}

	for _, resource := range request.Resources {
		if _, exists := quotaResourcesMap[resource.ResourceID]; !exists {
			return nil, 0, apiError.NewBadRequestError(fmt.Errorf("resource %s not found in organization quota", resource.ResourceID))
		}

		resourceQuantity := quotaResourcesMap[resource.ResourceID]

		if _, duplicate := seenReq[resource.ResourceID]; duplicate {
			return nil, 0, apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", resource.ResourceID))
		}
		seenReq[resource.ResourceID] = struct{}{}

		usage, err := u.ticketRepo.GetResourceUsage(request.NamespaceID, request.QuotaID, resource.ResourceID)
		if err != nil {
			return nil, 0, apiError.NewInternalServerError(err)
		}

		if resource.Quantity <= 0 {
			return nil, 0, apiError.NewForbiddenError(fmt.Errorf("request contain invalid resource quantity"))
		}

		maxQuota := resourceQuantity.Quantity
		if usage+resource.Quantity > maxQuota {
			return nil, 0, apiError.NewForbiddenError(fmt.Errorf("namespace usage exceeds quota limit for resource"))
		}

		// if request.Duration < 1800 {
		// 	return nil, 0, apiError.NewBadRequestError(fmt.Errorf("duration must be greater than 1800 seconds"))
		// }
		if request.Duration > resourceQuantity.ResourceProp.MaxDuration {
			return nil, 0, apiError.NewForbiddenError(fmt.Errorf("duration exceeds max limit for resource"))
		}

		price += resourceQuantity.ResourceProp.Price * float32(resource.Quantity)
		totalCredit += float32(resource.Quantity) * resourceQuantity.ResourceProp.Price * float32(request.Duration) / 3600
	}

	namespace, err := u.namespaceRepo.GetNamespaceByID(request.NamespaceID)
	if err != nil {
		return nil, 0, apiError.NewInternalServerError(err)
	}

	if totalCredit > namespace.Credit {
		return nil, 0, apiError.NewForbiddenError(fmt.Errorf("not enough credit want to use %.2f but only have %.2f", totalCredit, namespace.Credit))
	}
	namespace.Credit -= totalCredit

	return namespace, price, nil
}

func (u *TicketUsecase) createTicket(quota *models.NamespaceQuota, request *dtos.CreateTicketRequest, userID uuid.UUID, price float32) (*models.Ticket, error) {
	node, err := u.resourceRepo.GetResourceNodeByID(quota.NodeID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	ticket := &models.Ticket{
		NamespaceID:    request.NamespaceID,
		Name:           request.Name,
		Duration:       request.Duration,
		OwnerID:        userID,
		QuotaID:        request.QuotaID,
		NodeID:         quota.NodeID,
		ResourcePoolID: node.ResourcePoolID,
		GlideletURN:    node.ResourcePool.GlideletURN,
		Price:          price,
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

	ticket, err = u.ticketRepo.GetTicketByID(ticket.ID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return ticket, nil
}

func (u *TicketUsecase) CreateTicket(request *dtos.CreateTicketRequest, userID uuid.UUID) (*dtos.GliderTicketResponse, error) {
	if err := u.isNamespaceMember(request.NamespaceID, userID); err != nil {
		return nil, err
	}

	isAssigned, err := u.quotaRepo.IsAssigned(request.NamespaceID, request.QuotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !isAssigned {
		return nil, apiError.NewForbiddenError("quota is not assigned to the namespace")
	}

	quota, err := u.quotaRepo.GetNamespaceQuotaByID(request.QuotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	quotaResourcesMap := make(map[uuid.UUID]models.ResourceQuantity)
	for _, resource := range quota.Resources {
		quotaResourcesMap[resource.ResourceProp.ResourceID] = resource
	}

	namespace, ticketPrice, err := u.validateCreateTicketRequest(request, quotaResourcesMap)
	if err != nil {
		return nil, err
	}

	err = u.namespaceRepo.UpdateNamespace(namespace)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	ticket, err := u.createTicket(quota, request, userID, ticketPrice)
	if err != nil {
		return nil, err
	}

	return u.formatTicketResponse(ticket), nil
}

func (u *TicketUsecase) CancelTicket(ticketID uuid.UUID, userID uuid.UUID) error {
	ticket, err := u.ticketRepo.GetTicketByID(ticketID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	if err := u.isNamespaceMember(ticket.NamespaceID, userID); err != nil {
		return err
	}

	if ticket.OwnerID != userID {
		return apiError.NewForbiddenError("user is not the owner of the ticket")
	}

	if ticket.Status != "created" {
		return apiError.NewBadRequestError("only tickets in created status can be cancelled")
	}

	err = u.ticketRepo.CancelTicket(ticketID, time.Now())
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	namespace, err := u.namespaceRepo.GetNamespaceByID(ticket.NamespaceID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	namespace.Credit += float32(ticket.Duration) / 3600 * ticket.Price
	err = u.namespaceRepo.UpdateNamespace(namespace)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	return nil
}

func (u *TicketUsecase) DeleteTicket(ticketID uuid.UUID, userID uuid.UUID) error {
	ticket, err := u.ticketRepo.GetTicketByID(ticketID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	if err := u.isNamespaceMember(ticket.NamespaceID, userID); err != nil {
		return err
	}

	if ticket.OwnerID != userID {
		return apiError.NewForbiddenError("user is not the owner of the ticket")
	}

	err = u.ticketRepo.DeleteTicket(ticketID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	return nil
}

func (u *TicketUsecase) GetNamespaceTickets(namespaceID uuid.UUID, userID uuid.UUID) ([]models.Ticket, error) {
	if err := u.isNamespaceMember(namespaceID, userID); err != nil {
		return nil, err
	}

	tickets, err := u.ticketRepo.GetTicketsByNamespaceID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return tickets, nil
}

func (u *TicketUsecase) GetTicket(ticketID uuid.UUID, userID uuid.UUID) (*models.Ticket, error) {
	ticket, err := u.ticketRepo.GetTicketByID(ticketID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	if err := u.isNamespaceMember(ticket.NamespaceID, userID); err != nil {
		return nil, err
	}

	if ticket.OwnerID != userID {
		return nil, apiError.NewForbiddenError("user is not the owner of the ticket")
	}

	return ticket, nil
}

func (u *TicketUsecase) GetUserTickets(userID uuid.UUID) ([]dtos.TicketResponse, error) {
	tickets, err := u.ticketRepo.GetTicketsByUserID(userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	var responses []dtos.TicketResponse
	for _, ticket := range tickets {
		ticketResponse := dtos.TicketResponse{
			ID:             ticket.ID.String(),
			Name:           ticket.Name,
			Status:         ticket.Status,
			StartTime:      ticket.StartTime,
			EndTime:        ticket.EndTime,
			CancelTime:     ticket.CancelTime,
			Duration:       ticket.Duration,
			Price:          ticket.Price,
			OwnerID:        ticket.OwnerID.String(),
			NamespaceID:    ticket.NamespaceID.String(),
			NamespaceName:  ticket.Namespace.Name,
			ProjectID:      ticket.Namespace.ProjectID.String(),
			ProjectName:    ticket.Namespace.Project.Name,
			NodeID:         ticket.NodeID,
			GlideletURN:    ticket.GlideletURN,
			ResourcePoolID: ticket.ResourcePoolID.String(),
			QuotaID:        ticket.QuotaID.String(),
			RedeemTimeout:  ticket.RedeemTimeout,
			Resources:      ticket.Resources,
		}
		responses = append(responses, ticketResponse)
	}

	return responses, nil
}

func (u *TicketUsecase) StartTicket(request *dtos.StartTicketsRequest) ([]models.Ticket, error) {
	var tickets []models.Ticket
	for _, ticket := range request.Tickets {
		t, err := u.ticketRepo.GetTicketByID(ticket.TicketID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}
		if t.Status != "created" {
			return nil, apiError.NewBadRequestError("ticket is not in created status")
		}
		err = u.ticketRepo.StartTicket(ticket.TicketID, ticket.StartTime)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}

		t, err = u.ticketRepo.GetTicketByID(ticket.TicketID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}
		tickets = append(tickets, *t)
	}
	return tickets, nil
}

func (u *TicketUsecase) StopPendingTicket(request *dtos.StopPendingTicketsRequest) error {
	errorTickets := []uuid.UUID{}
	for _, ticketID := range request.Tickets {
		err := u.ticketRepo.StopPendingTicket(ticketID)
		if err != nil {
			errorTickets = append(errorTickets, ticketID)
			continue
		}

	}
	if len(errorTickets) > 0 {
		return apiError.NewInternalServerError(fmt.Errorf("some tickets could not be stopped"))
	}
	return nil
}

func (u *TicketUsecase) StopTicket(request *dtos.StopTicketsRequest) ([]models.Ticket, error) {
	var tickets []models.Ticket
	for _, ticketID := range request.Tickets {
		t, err := u.ticketRepo.GetTicketByID(ticketID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}

		if t.Status != "running" {
			return nil, apiError.NewBadRequestError("ticket is not in running status")
		}
		endTime := request.StopTime
		err = u.ticketRepo.StopTicket(ticketID, endTime)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}

		actualSeconds := endTime.Sub(*t.StartTime).Seconds()
		if actualSeconds < float64(t.Duration) {
			namespace, err := u.namespaceRepo.GetNamespaceByID(t.NamespaceID)
			if err != nil {
				return nil, apiError.NewInternalServerError(err)
			}

			namespace.Credit += float32(float64(t.Duration)-actualSeconds) / 3600 * t.Price
			err = u.namespaceRepo.UpdateNamespace(namespace)
			if err != nil {
				return nil, apiError.NewInternalServerError(err)
			}
		}

		t, err = u.ticketRepo.GetTicketByID(ticketID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}
		tickets = append(tickets, *t)
	}
	return tickets, nil
}

func (u *TicketUsecase) ResetTickets(ticketIDs []uuid.UUID) error {
	return u.ticketRepo.ResetTickets(ticketIDs)
}

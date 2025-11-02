package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/signature_helper"
	"github.com/google/uuid"
)

type TicketUsecase struct {
	namespaceRepo namespaceInterfaces.NamespaceRepository
	ticketRepo    interfaces.TicketRepository
	quotaRepo     quotaInterfaces.QuotaRepository
	userRepo      userInterfaces.UsersRepository
}

func NewTicketUsecase(namespaceRepo namespaceInterfaces.NamespaceRepository, ticketRepo interfaces.TicketRepository, quotaRepo quotaInterfaces.QuotaRepository, userRepo userInterfaces.UsersRepository) interfaces.TicketUsecase {
	return &TicketUsecase{
		namespaceRepo: namespaceRepo,
		ticketRepo:    ticketRepo,
		quotaRepo:     quotaRepo,
		userRepo:      userRepo,
	}
}

func (u *TicketUsecase) CreateTicket(request *dtos.CreateTicketRequest, userID uuid.UUID) (*dtos.GliderTicketResponse, error) {
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
	price := float32(0)

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

		if resource.Quantity <= 0 {
			return nil, apiError.NewForbiddenError("resource quota is zero")
		}

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
		price = resourceQuantity.ResourceProp.Price * float32(resource.Quantity)
		totalCredit += float32(resource.Quantity) * resourceQuantity.ResourceProp.Price * float32(request.Duration)
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
		Duration:       request.Duration,
		OwnerID:        userID,
		QuotaID:        request.QuotaID,
		ResourcePoolID: quota.ResourcePoolID,
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

	// reload ticket with resources
	ticket, err = u.ticketRepo.GetTicketByID(ticket.ID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	log.Printf("Created ticket: %+v", ticket)
	gliderTicket := u.FormatTicketResponse(ticket)
	log.Printf("Formatted glider ticket: %+v", gliderTicket)
	return gliderTicket, nil
}

func (u *TicketUsecase) GetNamespaceTickets(namespaceID uuid.UUID, userID uuid.UUID) ([]models.Ticket, error) {
	isMember, err := u.isNamespaceMember(userID, namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !isMember {
		return nil, apiError.NewUnauthorizedError("user is not a member of the namespace")
	}

	tickets, err := u.ticketRepo.GetTicketsByNamespaceID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return tickets, nil
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

func (u *TicketUsecase) StopTicket(request *dtos.StopTicketsRequest) ([]models.Ticket, error) {
	var tickets []models.Ticket
	for _, ticket := range request.Tickets {
		t, err := u.ticketRepo.GetTicketByID(ticket.TicketID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}
		if t.Status != "running" {
			return nil, apiError.NewBadRequestError("ticket is not in running status")
		}
		endTime := time.Now()
		err = u.ticketRepo.StopTicket(ticket.TicketID, endTime)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}

		actualSeconds := endTime.Sub(*t.StartTime).Seconds()
		if actualSeconds < float64(t.Duration) {
			namespace, err := u.namespaceRepo.GetNamespaceByID(t.NamespaceID)
			if err != nil {
				return nil, apiError.NewInternalServerError(err)
			}

			namespace.Credit += float32(float64(t.Duration)-actualSeconds) * t.Price
			err = u.namespaceRepo.UpdateNamespace(namespace)
			if err != nil {
				return nil, apiError.NewInternalServerError(err)
			}
		}

		t, err = u.ticketRepo.GetTicketByID(ticket.TicketID)
		if err != nil {
			return nil, apiError.NewInternalServerError(err)
		}
		tickets = append(tickets, *t)
	}
	return tickets, nil
}

func (u *TicketUsecase) isNamespaceMember(userID uuid.UUID, namespaceID uuid.UUID) (bool, error) {
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return false, err
	}
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	if !helper.ContainsUserID(namespace.Members, userID) && (*user.NamespaceID != namespace.ID) {
		return false, nil
	}

	return true, nil
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
			ResourcePoolID: ticket.ResourcePoolID.String(),
			QuotaID:        ticket.QuotaID.String(),
			RedeemTimeout:  ticket.RedeemTimeout,
			Resources:      ticket.Resources,
		}
		responses = append(responses, ticketResponse)
	}

	return responses, nil
}

func (u *TicketUsecase) GetTicket(ticketID uuid.UUID, userID uuid.UUID) (*models.Ticket, error) {
	ticket, err := u.ticketRepo.GetTicketByID(ticketID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	isMember, err := u.isNamespaceMember(userID, ticket.NamespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !isMember && ticket.OwnerID != userID {
		return nil, apiError.NewUnauthorizedError("user is not a member of the namespace or the owner of the ticket")
	}

	return ticket, nil
}

func (u *TicketUsecase) FormatTicketResponse(ticket *models.Ticket) *dtos.GliderTicketResponse {
	gliderTicket := dtos.GliderTicket{
		ID:                ticket.ID,
		NamespaceURN:      ticket.Namespace.ID.String(),
		GlideletURN:       ticket.ResourcePoolID.String(),
		Spec:              u.FormatGliderSpec(ticket),
		ReferenceTicketID: "",
		RedeemTimeout:     ticket.RedeemTimeout,
		Lease:             ticket.Duration,
		CreatedAt:         ticket.CreatedAt,
	}

	signature, err := signature_helper.SignTicket(gliderTicket)
	if err != nil {
		return nil
	}

	return &dtos.GliderTicketResponse{
		Ticket:    gliderTicket,
		Signature: signature,
	}
}

func (u *TicketUsecase) FormatGliderSpec(ticket *models.Ticket) dtos.GliderSpec {
	var resources []dtos.SpecResource
	for _, r := range ticket.Resources {
		resources = append(resources, dtos.SpecResource{
			ResourceID: r.ResourceID.String(),
			Name:       r.Resource.Name,
			Quantity:   r.Quantity,
			Unit:       r.Resource.ResourceType.Unit,
		})
	}

	return dtos.GliderSpec{
		Type:      dtos.ResourceUnitTypeCPU,
		PoolID:    ticket.ResourcePoolID.String(),
		Resources: resources,
	}
}

func (u *TicketUsecase) CancelTicket(ticketID uuid.UUID, userID uuid.UUID) error {
	ticket, err := u.ticketRepo.GetTicketByID(ticketID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	isMember, err := u.isNamespaceMember(userID, ticket.NamespaceID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}
	if !isMember && ticket.OwnerID != userID {
		return apiError.NewUnauthorizedError("user is not a member of the namespace or the owner of the ticket")
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

	namespace.Credit += float32(ticket.Duration) * ticket.Price
	err = u.namespaceRepo.UpdateNamespace(namespace)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	return nil
}

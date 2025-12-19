package usecase

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/ClearingHouse/pkg/signature_helper"
	"github.com/google/uuid"
)

func (u *TicketUsecase) formatTicketResponse(ticket *models.Ticket) *dtos.GliderTicketResponse {
	gliderTicket := dtos.GliderTicket{
		ID:                ticket.ID,
		NamespaceID:       ticket.NamespaceID,
		NamespaceName:     ticket.Namespace.Name,
		ProjectID:         *ticket.Namespace.ProjectID,
		ProjectName:       ticket.Namespace.Project.Name,
		NodeID:            ticket.NodeID,
		NodeName:          ticket.Node.Name,
		ResourcePoolID:    ticket.ResourcePoolID,
		ResourcePoolName:  ticket.ResourcePool.Name,
		GlideletURN:       ticket.GlideletURN,
		GlideletName:      ticket.ResourcePool.Name,
		OrganizationName:  ticket.ResourcePool.Organization.Name,
		Spec:              u.FormatGliderSpec(ticket),
		ReferenceTicketID: uuid.Nil,
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
		PoolID:    ticket.NodeID.String(),
		Resources: resources,
	}
}

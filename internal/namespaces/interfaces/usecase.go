package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dto"
	"github.com/google/uuid"
)

type NamespacesUsecase interface {
	GetAll(projectID uuid.UUID) ([]models.Namespace, error)
	GetByID(namespaceID uuid.UUID) (*models.Namespace, error)
	Create(projectID uuid.UUID, request *dto.CreateNamespaceRequest) (*models.Namespace, error)
	AddMembers(userID uuid.UUID, namespaceID uuid.UUID, memberIDs uuid.UUIDs) (*models.Namespace, error)
	ChangeQuota(userID uuid.UUID, namespaceID uuid.UUID, quotaID uuid.UUID) (*models.Namespace, error)
	RequestTicket(userID uuid.UUID, namespaceID uuid.UUID, request *dto.CreateTicketRequest) (*models.Ticket, error)
	TerminateTicket(userID uuid.UUID, namespaceID uuid.UUID, ticketID uuid.UUID) error
}

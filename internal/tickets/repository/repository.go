package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) interfaces.TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) GetTicketByID(ticketID uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("Resources").First(&ticket, "id = ?", ticketID).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) CreateTicket(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) CreateTicketResource(resource *models.TicketResource) error {
	return r.db.Create(resource).Error
}

func (r *TicketRepository) GetNamespaceUsage(namespaceID, quotaID, resourceID uuid.UUID) (uint, error) {
	var total uint

	err := r.db.Model(&models.TicketResource{}).
		Select("COALESCE(SUM(quantity), 0)").
		Joins("JOIN tickets ON tickets.id = ticket_resources.ticket_id").
		Where("tickets.namespace_id = ? AND ticket_resources.quota_id = ? AND ticket_resources.resource_id = ?", namespaceID, quotaID, resourceID).
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

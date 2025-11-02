package repository

import (
	"time"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	"github.com/ClearingHouse/pkg/enum"
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
	err := r.db.Preload("Resources.Resource.ResourceType").Preload("ResourcePool.Organization").Preload("Namespace.Project").First(&ticket, "id = ?", ticketID).Error
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

func (r *TicketRepository) GetResourceUsage(namespaceID, quotaID, resourceID uuid.UUID) (uint, error) {
	var total uint

	err := r.db.Debug().Model(&models.TicketResource{}).
		Select("COALESCE(SUM(quantity), 0)").
		Joins("JOIN tickets ON tickets.id = ticket_resources.ticket_id").
		Where("tickets.namespace_id = ? AND ticket_resources.resource_id = ? AND tickets.quota_id = ? AND tickets.status IN ?", namespaceID, resourceID, quotaID, enum.UsingStatuses).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *TicketRepository) GetTicketsByNamespaceID(namespaceID uuid.UUID) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Preload("Resources").Where("namespace_id = ?", namespaceID).Order("created_at DESC").Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) GetTicketsByUserID(userID uuid.UUID) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Preload("Resources").Preload("Namespace.Project").Where("owner_id = ?", userID).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) StartTicket(ticketID uuid.UUID, startTime time.Time) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", ticketID).Updates(map[string]interface{}{
		"status":     "running",
		"start_time": startTime,
	}).Error

}

func (r *TicketRepository) StopTicket(ticketID uuid.UUID, stopTime time.Time) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", ticketID).Updates(map[string]interface{}{
		"status":   "stopped",
		"end_time": stopTime,
	}).Error
}

func (r *TicketRepository) CancelTicket(ticketID uuid.UUID, cancelTime time.Time) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", ticketID).Updates(map[string]interface{}{
		"status":      "cancelled",
		"cancel_time": cancelTime,
	}).Error
}

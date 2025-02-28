package repository

import (
	"github.com/ClearingHouse/internal/tickets/interfaces"
	"gorm.io/gorm"
)

type TicketsRepository struct {
	db *gorm.DB
}

func NewTicketsRepository(db *gorm.DB) interfaces.TicketsRepository {
	return &TicketsRepository{db: db}
}

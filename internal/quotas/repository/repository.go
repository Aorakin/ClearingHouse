package repository

import (
	"github.com/ClearingHouse/internal/quotas/interfaces"
	"gorm.io/gorm"
)

type QuotasRepository struct {
	db *gorm.DB
}

func NewQuotasRepository(db *gorm.DB) interfaces.QuotasRepository {
	return &QuotasRepository{db: db}
}

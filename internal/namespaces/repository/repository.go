package repository

import (
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	"gorm.io/gorm"
)

type NamespacesRepository struct {
	db *gorm.DB
}

func NewNamespacesRepository(db *gorm.DB) interfaces.NamespacesRepository {
	return &NamespacesRepository{db: db}
}

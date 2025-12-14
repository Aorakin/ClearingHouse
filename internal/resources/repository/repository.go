package repository

import (
	"github.com/ClearingHouse/internal/resources/interfaces"
	"gorm.io/gorm"
)

func NewResourceRepository(db *gorm.DB) (interfaces.ResourcePoolRepository, interfaces.ResourceRepository, interfaces.ResourceTypeRepository) {
	return &ResourcePoolRepository{db: db}, &ResourceRepository{db: db}, &ResourceTypeRepository{db: db}
}

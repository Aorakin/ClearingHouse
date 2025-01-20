package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type UsersRepository interface {
	GetUserGoogle(*oauth2.Token) (map[string]interface{}, error)
	Create(*models.User) error
	Delete(uuid.UUID) error
	GetUser(uuid.UUID) (*models.User, error)
	GetByUsername(string) (*models.User, error)
}

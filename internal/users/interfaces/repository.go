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
	GetByID(uuid.UUID) (*models.User, error)
	GetByUsername(string) (*models.User, error)
	FindOrCreateUser(email string, firstName string, lastName string) (*models.User, error)
	GetByIDs(userIDs []uuid.UUID) ([]models.User, error)
}

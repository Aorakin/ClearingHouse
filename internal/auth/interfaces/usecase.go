package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthUsecase interface {
	GenerateGoogleLoginURL(state string) string
	HandleGoogleCallback(string, *gin.Context) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
}

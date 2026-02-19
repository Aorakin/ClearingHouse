package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthUsecase interface {
	GenerateGoogleLoginURL(state string) string
	GenerateGoogleRegisterURL(state string) string
	HandleGoogleCallback(string, *gin.Context) (*models.User, error)
	HandleGoogleRegisterCallback(string, *gin.Context) (*models.User, error)
	ManualRegister(email, firstName, lastName string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GenerateTokens(user *models.User) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(refreshToken string) (string, error)
	BlacklistToken(token string) error
}

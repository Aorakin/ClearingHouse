package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthUsecase interface {
	GenerateGoogleLoginURL(state, portal string) string
	GenerateGoogleRegisterURL(state, portal string) string
	HandleGoogleCallback(code, portal string, c *gin.Context) (*models.User, error)
	HandleGoogleRegisterCallback(code, portal string, c *gin.Context) (*models.User, error)
	ManualRegister(email, firstName, lastName string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GenerateTokens(user *models.User) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(refreshToken string) (string, error)
	BlacklistToken(token string) error
}

package usecase

import (
	"context"

	"errors"

	"github.com/ClearingHouse/config"
	"github.com/ClearingHouse/internal/auth/interfaces"
	"github.com/ClearingHouse/internal/models"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type AuthUsecase struct {
	userRepo userInterfaces.UsersRepository
}

func NewAuthUsecase(userRepo userInterfaces.UsersRepository) interfaces.AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
	}
}

func (u *AuthUsecase) GenerateGoogleLoginURL(state string) string {
	return config.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (u *AuthUsecase) HandleGoogleCallback(code string, c *gin.Context) (*models.User, error) {
	token, err := config.GoogleOauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	googleUser, err := u.userRepo.GetUserGoogle(token)
	if err != nil {
		return nil, err
	}

	email, emailOk := googleUser["email"].(string)
	firstName, firstNameOk := googleUser["given_name"].(string)
	lastName, lastNameOk := googleUser["family_name"].(string)
	if !emailOk || !firstNameOk || !lastNameOk {
		return nil, errors.New("invalid user data from Google")
	}

	user, err := u.userRepo.FindOrCreateUser(email, firstName, lastName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return u.userRepo.GetByID(userID)
}

// func (u *AuthUsecase) GenerateToken(email string) (string, error) {
// 	// Generate a JWT token for the user
// 	token, err := createJWT(email)
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }

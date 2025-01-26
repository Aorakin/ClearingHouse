package usecase

import (
	"context"
	"errors"

	"github.com/ClearingHouse/config"
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/users/dtos"
	"github.com/ClearingHouse/internal/users/interfaces"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type UsersUsecase struct {
	usersRepository interfaces.UsersRepository
}

func NewUsersUsecase(userRepository interfaces.UsersRepository) interfaces.UsersUsecase {
	return &UsersUsecase{
		usersRepository: userRepository,
	}
}

func (u *UsersUsecase) GenerateLoginURL(state string) string {
	return config.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (u *UsersUsecase) HandleGoogleCallback(code string, c *gin.Context) (map[string]interface{}, error) {
	token, err := config.GoogleOauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}
	return u.usersRepository.GetUserGoogle(token)
}

func (u *UsersUsecase) Register(registerInput dtos.RegisterInput) error {
	if _, err := u.usersRepository.GetByUsername(registerInput.Username); err == nil {
		return errors.New("username already in used")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInput.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("invalid username or password")
	}

	user := &models.User{
		Username: registerInput.Username,
		Password: string(hashedPassword),
	}

	return u.usersRepository.Create(user)
}

func (u *UsersUsecase) Login(username, password string) (*models.User, error) {
	user, err := u.usersRepository.GetByUsername(username)
	if err != nil || user == nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

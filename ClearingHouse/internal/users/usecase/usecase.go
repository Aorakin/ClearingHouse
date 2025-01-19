package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ClearingHouse/config"
	"github.com/ClearingHouse/internal/users/interfaces"
	"github.com/gin-gonic/gin"
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
	c.SetCookie("access_token", token.AccessToken, int(token.Expiry.Sub(time.Now()).Seconds()), "/", "localhost", true, true)
	fmt.Println("hello")
	return u.usersRepository.GetUserGoogle(token)
}

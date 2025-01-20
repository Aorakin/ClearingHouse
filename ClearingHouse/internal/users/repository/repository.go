package repository

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/users/interfaces"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) interfaces.UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetUserGoogle(token *oauth2.Token) (map[string]interface{}, error) {
	client := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info")
	}

	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (r *UsersRepository) Create(user *models.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) GetUser(id uuid.UUID) (*models.User, error) {
	var user *models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UsersRepository) GetByUsername(username string) (*models.User, error) {
	var user *models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return user, nil
}

package helper

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func ContainsUserID(users []models.User, userID uuid.UUID) bool {
	for _, u := range users {
		if u.ID == userID {
			return true
		}
	}
	return false
}

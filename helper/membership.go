package helper

import (
	"log"

	"github.com/ClearingHouse/internal/models"
	"github.com/google/uuid"
)

func ContainsUserID(users []models.User, userID uuid.UUID) bool {
	log.Printf("Checking if userID %s is in users list", userID)
	for _, u := range users {
		log.Printf("Comparing with user ID: %s", u.ID)
		if u.ID == userID {
			return true
		}
	}
	return false
}

package main

import (
	"log"
	"time"

	"github.com/ClearingHouse/internal/app"
	AuthRepository "github.com/ClearingHouse/internal/auth/repository"
	"gorm.io/gorm"
)

func main() {
	db, err := app.InitDataBase()
	if err != nil {
		log.Fatal(err)
	}
	defer app.CloseDB()

	// Start token blacklist cleanup job
	go startTokenCleanupJob(db)

	app := app.NewApp(db)
	if err := app.Run(); err != nil {
		log.Fatal("Server encountered an error:", err)
	}
}

// startTokenCleanupJob runs a periodic cleanup of expired tokens from the blacklist
func startTokenCleanupJob(db *gorm.DB) {
	blacklistRepo := AuthRepository.NewTokenBlacklistRepository(db)
	ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	log.Println("Running initial token blacklist cleanup...")
	if err := blacklistRepo.CleanupExpired(); err != nil {
		log.Printf("Token cleanup failed: %v", err)
	} else {
		log.Println("Token blacklist cleanup completed successfully")
	}

	for range ticker.C {
		log.Println("Running scheduled token blacklist cleanup...")
		if err := blacklistRepo.CleanupExpired(); err != nil {
			log.Printf("Token cleanup failed: %v", err)
		} else {
			log.Println("Token blacklist cleanup completed successfully")
		}
	}
}

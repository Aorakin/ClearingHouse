package main

import (
	"log"
	"time"

	"github.com/ClearingHouse/internal/app"
	AuthRepository "github.com/ClearingHouse/internal/auth/repository"
)

// This script cleans up expired tokens from the blacklist
// Run this as a cron job or scheduled task
func main() {
	db, err := app.InitDataBase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer app.CloseDB()

	blacklistRepo := AuthRepository.NewTokenBlacklistRepository(db)

	log.Println("Starting token blacklist cleanup...")
	err = blacklistRepo.CleanupExpired()
	if err != nil {
		log.Fatal("Failed to cleanup expired tokens:", err)
	}

	log.Println("Token blacklist cleanup completed successfully at", time.Now())
}

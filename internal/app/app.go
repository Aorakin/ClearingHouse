package app

import (
	"fmt"
	"log"

	"github.com/ClearingHouse/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	gin        *gin.Engine
	postgresDB *gorm.DB
	// config      *config.Config
}

func NewApp(postgresDB *gorm.DB) *App {
	return &App{
		gin:        gin.New(),
		postgresDB: postgresDB,
	}
}

func (s *App) Run() error {
	err := s.gin.SetTrustedProxies([]string{"192.168.0.0/16", "10.0.0.0/8"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	config.InitConfig()
	store := config.NewSessionStore("ClearingHouseSession", 3600)
	s.gin.RouterGroup.Use(sessions.Sessions("ClearingHouseSession", store))

	if err := s.MapHandlers(); err != nil {
		return err
	}

	serverURL := fmt.Sprintf(":%s", "8080")
	return s.gin.Run(serverURL)
}

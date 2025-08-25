package app

import (
	"fmt"
	"log"
	"time"

	"github.com/ClearingHouse/config"
	"github.com/ClearingHouse/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	s.gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	config.InitConfig()
	store := config.NewSessionStore("session", 1800) // 30 min
	s.gin.Use(sessions.Sessions("session", store))

	docs.SwaggerInfo.Title = "ClearingHouse API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	if err := s.MapHandlers(); err != nil {
		return err
	}

	// Serve Swagger UI
	s.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	serverURL := fmt.Sprintf(":%s", "8080")
	return s.gin.Run(serverURL)
}

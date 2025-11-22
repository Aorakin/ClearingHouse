package app

import (
	"fmt"
	"log"
	"os"

	"github.com/ClearingHouse/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDataBase() (*gorm.DB, error) {
	err := godotenv.Load()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}
	DB.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.Project{},
		&models.Namespace{},
		&models.ResourcePool{},
		&models.ResourceType{},
		&models.Resource{},
		&models.ResourceProperty{},
		&models.OrganizationQuota{},
		&models.ProjectQuota{},
		&models.NamespaceQuota{},
		&models.ResourceQuantity{},
		&models.ResourceProperty{},
		&models.Ticket{},
		&models.TicketResource{},
	)
	return DB, nil
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting database connection: %v", err)
	}
	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("error closing the database connection: %v", err)
	}

	log.Println("Database connection closed.")
	return nil
}

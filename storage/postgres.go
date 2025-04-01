package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDB *gorm.DB

// InitPostgres initializes PostgreSQL connection
func InitPostgres() *gorm.DB {
	_ = godotenv.Load() // Load .env file

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Transaction{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
	PostgresDB = db
	return db
}

package db

import (
	"fmt"
	"log"
	"time"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDB *gorm.DB

// InitPostgresDB initializes the PostgreSQL database connection with retry mechanism
func InitPostgresDB() {
	dsn := "user=postgres password=postgres dbname=banking_db host=postgres_db port=5432 sslmode=disable"
	var err error

	// Retry logic
	for i := 0; i < 10; i++ {
		PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Postgres connected successfully!")
			return
		}
		log.Printf("Failed to connect to PostgreSQL (attempt %d/10): %v", i+1, err)
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Fatalf("Failed to connect to PostgreSQL after 10 attempts: %v", err)
}

// AutoMigrate ensures the account table is created in PostgreSQL
func AutoMigrate() {
	err := PostgresDB.AutoMigrate(&models.Account{})
	if err != nil {
		log.Fatalf("Failed to migrate Postgres DB: %v", err)
	}
}

package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

// InitMongo initializes MongoDB connection
func InitMongo() *mongo.Database {
	_ = godotenv.Load() // Load .env file

	mongoURI := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?authSource=%s",
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_AUTH_SOURCE"),
	)

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	MongoDB = client.Database(os.Getenv("DB_NAME")) // Use correct DB name
	log.Println("Connected to MongoDB successfully")
	return MongoDB
}

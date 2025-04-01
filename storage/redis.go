package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// InitRedis initializes the Redis client and returns it
func InitRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_HOST")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default if not set
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // No password set
		DB:       0,  // Default DB
	})

	// Ping Redis to check the connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
	return client
}

// SetSession stores a session token in Redis
func SetSession(client *redis.Client, token string, username string) error {
	return client.Set(context.Background(), token, username, 24*time.Hour).Err()
}

// GetSession retrieves a session from Redis
func GetSession(client *redis.Client, token string) (string, error) {
	return client.Get(context.Background(), token).Result()
}

// DeleteSession removes a session (Logout)
func DeleteSession(client *redis.Client, token string) error {
	return client.Del(context.Background(), token).Err()
}

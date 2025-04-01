package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ashil-poojary/banking-ledger-service/api/routes"
	"github.com/ashil-poojary/banking-ledger-service/config"
	"github.com/ashil-poojary/banking-ledger-service/storage"
	"github.com/ashil-poojary/banking-ledger-service/worker"
	"github.com/gorilla/mux"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database connections
	postgresDB := storage.InitPostgres()
	mongoDB := storage.InitMongo()
	redisClient := storage.InitRedis()

	// Initialize RabbitMQ
	_, rabbitMQChannel := storage.InitRabbitMQ()
	defer storage.CloseRabbitMQ()

	// Start transaction worker in a goroutine
	go worker.ProcessTransactions("transactions", postgresDB, mongoDB, rabbitMQChannel)

	// Start API server
	r := mux.NewRouter()
	routes.SetupRoutes(r, postgresDB, mongoDB, redisClient, rabbitMQChannel)

	fmt.Println("API server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

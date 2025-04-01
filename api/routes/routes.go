package routes

import (
	"github.com/ashil-poojary/banking-ledger-service/api/handlers"
	"github.com/ashil-poojary/banking-ledger-service/api/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// SetupRoutes initializes API routes
func SetupRoutes(r *mux.Router, postgresDB *gorm.DB, mongoDB *mongo.Database, redisClient *redis.Client, rabbitMQChannel *amqp.Channel) {
	accountHandler := handlers.NewAccountHandler(postgresDB)
	transactionHandler := handlers.NewTransactionHandler(postgresDB, mongoDB, rabbitMQChannel, "transactions")
	authHandler := handlers.NewAuthHandler(postgresDB, redisClient)

	r.Use(middleware.LoggingMiddleware)

	// Auth Routes
	r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/logout", authHandler.Logout).Methods("POST")

	// Secure account & transaction routes with middleware
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(redisClient))

	// Account Routes
	protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts/user", accountHandler.GetUserAccounts).Methods("GET")
	protected.HandleFunc("/accounts/{id}", accountHandler.GetAccount).Methods("GET")
	protected.HandleFunc("/accounts/{id}", accountHandler.UpdateAccount).Methods("PUT")
	protected.HandleFunc("/accounts/{id}", accountHandler.DeleteAccount).Methods("DELETE")

	// Transaction Routes
	protected.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods("POST")
	protected.HandleFunc("/transactions/{id}", transactionHandler.GetTransaction).Methods("GET")
	protected.HandleFunc("/transactions/history", transactionHandler.GetTransactionHistory).Methods("GET")
}

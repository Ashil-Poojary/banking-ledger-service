package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ashil-poojary/banking-ledger-service/internal/db"
	"github.com/ashil-poojary/banking-ledger-service/internal/queue"
	"github.com/ashil-poojary/banking-ledger-service/services"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Banking Ledger Service...")

	// Initialize PostgreSQL
	log.Println("Initializing PostgreSQL...")
	db.InitPostgresDB()
	log.Println("PostgreSQL initialized successfully.")

	// Initialize MongoDB
	log.Println("Initializing MongoDB...")
	db.InitMongoDB()
	log.Println("MongoDB initialized successfully.")

	// AutoMigrate to ensure tables exist
	log.Println("Running database migrations...")
	db.AutoMigrate()
	log.Println("Database migrations completed.")

	// Start processing transaction queue
	log.Println("Starting transaction queue processing...")
	go queue.ProcessTransactionQueue()

	// Initialize router
	r := mux.NewRouter()

	// Create account endpoint
	r.HandleFunc("/api/accounts", logRequest(services.CreateAccountHandler)).Methods("POST")
	r.HandleFunc("/api/transactions", logRequest(services.DepositHandler)).Methods("POST")
	r.HandleFunc("/api/transactions/{accountID}", logRequest(services.GetTransactionsHandler)).Methods("GET")
	r.HandleFunc("/api/transactions/withdraw", logRequest(services.WithdrawHandler)).Methods("POST")

	// Start the server
	log.Println("Banking Ledger Service is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Middleware to log incoming requests
func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, start.Format(time.RFC3339))
		handler(w, r)
		duration := time.Since(start)
		log.Printf("[%s] Completed in %v", r.Method, duration)
	}
}

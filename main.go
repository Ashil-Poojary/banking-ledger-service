package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// // Routes
	// r.HandleFunc("/accounts", CreateAccount).Methods("POST")
	// r.HandleFunc("/accounts/{id}", GetAccount).Methods("GET")
	// r.HandleFunc("/transactions", CreateTransaction).Methods("POST")
	// r.HandleFunc("/transactions/{id}", GetTransactions).Methods("GET")

	log.Println("API Gateway running on port 8080...")
	http.ListenAndServe(":8080", r)
}

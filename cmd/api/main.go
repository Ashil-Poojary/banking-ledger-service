package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Sample routes (implement actual handlers later)
	r.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Account Created"))
	}).Methods("POST")

	r.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Transaction Processed"))
	}).Methods("POST")

	log.Println("API Gateway running on port 8080...")
	http.ListenAndServe(":8080", r)
}

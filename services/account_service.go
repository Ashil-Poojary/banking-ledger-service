package services

import (
	"encoding/json"
	"net/http"

	"github.com/ashil-poojary/banking-ledger-service/internal/db"
	"github.com/ashil-poojary/banking-ledger-service/models"
)

// Create Account
func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if result := db.PostgresDB.Create(&account); result.Error != nil {
		http.Error(w, "Error creating account: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// Deposit Money
func DepositHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	transaction.Type = "deposit"
	if result := db.PostgresDB.Create(&transaction); result.Error != nil {
		http.Error(w, "Error processing deposit: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

// Withdraw Money
func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	transaction.Type = "withdraw"
	if result := db.PostgresDB.Create(&transaction); result.Error != nil {
		http.Error(w, "Error processing withdrawal: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

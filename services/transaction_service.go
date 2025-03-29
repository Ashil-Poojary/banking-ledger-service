package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ashil-poojary/banking-ledger-service/internal/db"
	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the account ID from URL parameters
	vars := mux.Vars(r)
	accountID := vars["accountID"]

	// Retrieve transactions from MongoDB
	collection := db.MongoDB.Database("banking_db").Collection("transactions")
	cursor, err := collection.Find(r.Context(), bson.M{"account_id": accountID})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving transactions: %v", err), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var transactions []models.Transaction
	for cursor.Next(r.Context()) {
		var transaction models.Transaction
		err := cursor.Decode(&transaction)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding transaction: %v", err), http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error during cursor iteration: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the transactions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

// func DepositHandler(w http.ResponseWriter, r *http.Request) {
// 	var transaction models.Transaction
// 	err := json.NewDecoder(r.Body).Decode(&transaction)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Process deposit into account using GORM
// 	var account models.Account
// 	if err := db.PostgresDB.First(&account, transaction.AccountID).Error; err != nil {
// 		http.Error(w, fmt.Sprintf("Account not found: %v", err), http.StatusNotFound)
// 		return
// 	}

// 	// Update balance in the account
// 	account.Balance += transaction.Amount
// 	if err := db.PostgresDB.Save(&account).Error; err != nil {
// 		http.Error(w, fmt.Sprintf("Error updating account balance: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Log the transaction in MongoDB
// 	collection := db.MongoDB.Database("banking_db").Collection("transactions")
// 	_, err = collection.InsertOne(r.Context(), transaction)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error logging transaction: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(transaction)
// }

// func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
// 	var transaction models.Transaction
// 	err := json.NewDecoder(r.Body).Decode(&transaction)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Process withdrawal from account using GORM
// 	var account models.Account
// 	if err := db.PostgresDB.First(&account, transaction.AccountID).Error; err != nil {
// 		http.Error(w, fmt.Sprintf("Account not found: %v", err), http.StatusNotFound)
// 		return
// 	}

// 	// Check if sufficient balance is available
// 	if account.Balance < transaction.Amount {
// 		http.Error(w, "Insufficient funds", http.StatusBadRequest)
// 		return
// 	}

// 	// Update balance in the account
// 	account.Balance -= transaction.Amount
// 	if err := db.PostgresDB.Save(&account).Error; err != nil {
// 		http.Error(w, fmt.Sprintf("Error updating account balance: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Log the transaction in MongoDB
// 	collection := db.MongoDB.Database("banking_db").Collection("transactions")
// 	_, err = collection.InsertOne(r.Context(), transaction)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error logging transaction: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(transaction)
// }

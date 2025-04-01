package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/ashil-poojary/banking-ledger-service/utils"
	"github.com/ashil-poojary/banking-ledger-service/worker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// TransactionHandler handles transaction-related requests
type TransactionHandler struct {
	PostgresDB *gorm.DB
	MongoDB    *mongo.Database
	RabbitMQ   worker.RabbitMQPublisher // ✅ Use the interface
	QueueName  string
}

// NewTransactionHandler initializes a new TransactionHandler
func NewTransactionHandler(postgresDB *gorm.DB, mongoDB *mongo.Database, rabbitMQ worker.RabbitMQPublisher, queueName string) *TransactionHandler {
	return &TransactionHandler{PostgresDB: postgresDB, MongoDB: mongoDB, RabbitMQ: rabbitMQ, QueueName: queueName}
}

// TransferFunds handles money transfers between accounts
func (h *TransactionHandler) TransferFunds(w http.ResponseWriter, r *http.Request) {
	var transferReq models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request payload")
		return
	}

	// Validate request data
	if transferReq.SourceAccount == "" || transferReq.DestinationAccount == "" || transferReq.Amount <= 0 {
		utils.SendResponse(w, http.StatusBadRequest, false, "", nil, "Invalid transfer details")
		return
	}

	// Start a transaction in PostgreSQL
	tx := h.PostgresDB.Begin()

	// Fetch source account
	var sourceAccount models.Account
	if err := tx.Where("account_number = ?", transferReq.SourceAccount).First(&sourceAccount).Error; err != nil {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, false, "", nil, "Source account not found")
		return
	}

	// Check if the source account has enough balance
	if sourceAccount.Balance < transferReq.Amount {
		tx.Rollback()
		utils.SendResponse(w, http.StatusBadRequest, false, "", nil, "Insufficient funds")
		return
	}

	// Fetch destination account
	var destinationAccount models.Account
	if err := tx.Where("account_number = ?", transferReq.DestinationAccount).First(&destinationAccount).Error; err != nil {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, false, "", nil, "Destination account not found")
		return
	}

	// Deduct from source and add to destination
	sourceAccount.Balance -= transferReq.Amount
	destinationAccount.Balance += transferReq.Amount

	if err := tx.Save(&sourceAccount).Error; err != nil || tx.Save(&destinationAccount).Error != nil {
		tx.Rollback()
		utils.SendResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to process transfer")
		return
	}

	// Commit transaction
	tx.Commit()

	// Create transaction log
	txn := models.Transaction{
		SourceAccount:      transferReq.SourceAccount,
		DestinationAccount: transferReq.DestinationAccount,
		AccountNumber:      transferReq.SourceAccount,
		Amount:             transferReq.Amount,
		Currency:           transferReq.Currency,
		Type:               "transfer",
		Status:             "completed",
	}

	// Store in MongoDB
	_, err := h.MongoDB.Collection("transactions").InsertOne(context.TODO(), txn)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to log transaction")
		return
	}

	// Publish event to RabbitMQ
	if err := worker.PublishTransaction(txn, h.RabbitMQ, h.QueueName); err != nil {
		utils.SendResponse(w, http.StatusServiceUnavailable, false, "", nil, "RabbitMQ is unavailable. Please try again later.")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Transfer successful", txn, "")
}

// GetTransaction retrieves a specific transaction from PostgreSQL
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.URL.Query().Get("account_number") // 🔹 Ensure correct query param name

	if accountNumber == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, "", nil, "Account number is required")
		return
	}

	var transactions []models.Transaction
	err := h.PostgresDB.Where("account_number = ?", accountNumber).Find(&transactions).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to retrieve transactions")
		return
	}

	if len(transactions) == 0 {
		utils.SendResponse(w, http.StatusNotFound, false, "", nil, "No transactions found for this account")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Transactions retrieved successfully", transactions, "")
}

// GetTransactionHistory retrieves transaction logs from MongoDB
func (h *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.URL.Query().Get("account_number") // 🔹 Ensure correct query param name

	var filter bson.M
	if accountNumber != "" {
		filter = bson.M{"account_number": accountNumber} // 🔹 Filter by account if provided
	} else {
		filter = bson.M{} // 🔹 Retrieve all transactions if no filter is given
	}

	cursor, err := h.MongoDB.Collection("transactions").Find(context.TODO(), filter)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to retrieve transaction history")
		return
	}
	defer cursor.Close(context.TODO())

	var transactions []models.Transaction
	if err = cursor.All(context.TODO(), &transactions); err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, "", nil, "Error decoding transaction history")
		return
	}

	log.Println("[DEBUG] Transaction history:", transactions)
	utils.SendResponse(w, http.StatusOK, true, "Transaction history retrieved successfully", transactions, "")
}

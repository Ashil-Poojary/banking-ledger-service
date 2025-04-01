package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/ashil-poojary/banking-ledger-service/utils"
	"github.com/ashil-poojary/banking-ledger-service/worker"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// TransactionHandler handles transaction-related requests
type TransactionHandler struct {
	PostgresDB *gorm.DB
	MongoDB    *mongo.Database
	RabbitMQ   *amqp.Channel
	QueueName  string
}

// NewTransactionHandler initializes a new TransactionHandler
func NewTransactionHandler(postgresDB *gorm.DB, mongoDB *mongo.Database, rabbitMQ *amqp.Channel, queueName string) *TransactionHandler {
	return &TransactionHandler{PostgresDB: postgresDB, MongoDB: mongoDB, RabbitMQ: rabbitMQ, QueueName: queueName}
}

// CreateTransaction handles transaction creation
// CreateTransaction handles transaction requests
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var txn models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request payload")
		return
	}

	// Publish transaction using RabbitMQ
	if err := worker.PublishTransaction(txn, h.RabbitMQ, h.QueueName); err != nil {
		utils.SendResponse(w, http.StatusServiceUnavailable, false, "", nil, "RabbitMQ is unavailable. Please try again later.")
		return
	}

	utils.SendResponse(w, http.StatusAccepted, true, "Transaction queued for processing", nil, "")
}

// GetTransaction retrieves a specific transaction from PostgreSQL
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, "", nil, "Invalid transaction ID")
		return
	}

	var txn models.Transaction
	if err := h.PostgresDB.First(&txn, id).Error; err != nil {
		utils.SendResponse(w, http.StatusNotFound, false, "", nil, "Transaction not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, "Transaction retrieved successfully", txn, "")
}

// GetTransactionHistory retrieves transaction logs from MongoDB
func (h *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.URL.Query().Get("accountNumber")
	if accountNumber == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, "", nil, "Account number is required")
		return
	}

	filter := bson.M{"accountNumber": accountNumber}
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

	utils.SendResponse(w, http.StatusOK, true, "Transaction history retrieved successfully", transactions, "")
}

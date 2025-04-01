package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// ProcessTransactions listens for transaction messages and processes them
func ProcessTransactions(queueName string, postgresDB *gorm.DB, mongoDB *mongo.Database, rabbitMQChannel *amqp.Channel) {
	log.Println("[Worker] Starting transaction processing...")

	// Ensure the queue exists before consuming
	_, err := rabbitMQChannel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("[Worker] Failed to declare queue: %v", err)
	}

	msgs, err := rabbitMQChannel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Worker] Failed to register consumer: %v", err)
	}

	for msg := range msgs {
		var transaction models.Transaction

		if err := json.Unmarshal(msg.Body, &transaction); err != nil {
			log.Println("[Worker] Failed to unmarshal transaction:", err)
			msg.Reject(false) // Permanently reject invalid messages
			continue
		}

		// 🔹 **PostgreSQL Transaction**
		err := processTransaction(postgresDB, mongoDB, &transaction)
		if err != nil {
			log.Println("[Worker] Transaction processing failed:", err)
			msg.Nack(false, true) // Retry message in RabbitMQ
			continue
		}

		msg.Ack(false) // Acknowledge successful processing
	}
}

// processTransaction handles transaction logic within a PostgreSQL transaction
func processTransaction(postgresDB *gorm.DB, mongoDB *mongo.Database, transaction *models.Transaction) error {
	tx := postgresDB.Begin() // Start transaction
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("[Worker] Panic recovered:", r)
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	var balance float64

	// 🔹 **Lock the row for update to prevent race conditions**
	err := tx.Raw(`SELECT balance FROM accounts WHERE account_number = ? FOR UPDATE`, transaction.AccountNumber).Scan(&balance).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 🔹 **Check sufficient funds for withdrawal**
	if transaction.Type == "withdrawal" && balance < transaction.Amount {
		tx.Rollback()
		return errors.New("insufficient funds for withdrawal")
	}

	// 🔹 **Update balance**
	newBalance := balance
	if transaction.Type == "withdrawal" {
		newBalance -= transaction.Amount
	} else if transaction.Type == "deposit" {
		newBalance += transaction.Amount
	}

	err = tx.Exec(`UPDATE accounts SET balance = ? WHERE account_number = ?`, newBalance, transaction.AccountNumber).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 🔹 **Commit PostgreSQL transaction**
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 🔹 **Insert transaction into MongoDB with retry**
	transaction.Status = "completed"
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		_, err := mongoDB.Collection("transactions").InsertOne(context.TODO(), transaction)
		if err == nil {
			log.Println("[Worker] Transaction successfully logged in MongoDB")
			return nil // Success
		}
		log.Printf("[Worker] MongoDB insertion failed (attempt %d/%d). Retrying in 2s... Error: %v", i+1, retryCount, err)
		time.Sleep(2 * time.Second) // Backoff before retry
	}

	return errors.New("failed to insert transaction log into MongoDB after retries")
}

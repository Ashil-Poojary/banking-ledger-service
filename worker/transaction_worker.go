package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// ProcessTransactions listens for transaction messages and processes them
func ProcessTransactions(queueName string, postgresDB *gorm.DB, mongoDB *mongo.Database, rabbitMQChannel *amqp.Channel) {
	log.Println("[Worker] Starting transaction processing...")

	q, err := rabbitMQChannel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("[Worker] Failed to declare queue: %v", err)
	}
	log.Printf("[Worker] Listening on queue: %s", queueName)

	msgs, err := rabbitMQChannel.Consume(
		q.Name,
		"",
		false, // Manual acknowledgment for reliability
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[Worker] Failed to register consumer: %v", err)
	}

	log.Println("[Worker] Waiting for messages...")

	for msg := range msgs {
		log.Println("[Worker] Received new transaction message")
		var transaction models.Transaction
		if err := json.Unmarshal(msg.Body, &transaction); err != nil {
			log.Println("[Worker] Failed to unmarshal transaction:", err)
			msg.Reject(false) // Permanent failure, discard message
			continue
		}

		log.Printf("[Worker] Processing transaction: %+v", transaction)

		// Assign a new MongoDB ID
		transaction.ID = primitive.NewObjectID()
		transaction.Status = "completed"

		// Start a PostgreSQL transaction
		tx := postgresDB.Begin()
		if tx.Error != nil {
			log.Println("[Worker] Failed to start transaction:", tx.Error)
			msg.Nack(false, true) // Requeue message
			continue
		}
		log.Println("[Worker] Started PostgreSQL transaction")

		// Ensure account exists before updating
		var account models.Account
		if err := tx.Where("account_number = ?", transaction.AccountNumber).First(&account).Error; err != nil {
			log.Printf("[Worker] Account not found (%s), rejecting transaction", transaction.AccountNumber)
			tx.Rollback()
			msg.Reject(false) // Permanent failure, discard message
			continue
		}
		log.Printf("[Worker] Account found: %+v", account)

		// Update account balance
		result := tx.Exec("UPDATE accounts SET balance = balance + ? WHERE account_number = ?", transaction.Amount, transaction.AccountNumber)
		if result.Error != nil {
			log.Println("[Worker] Failed to update balance:", result.Error)
			tx.Rollback()
			msg.Nack(false, true) // Requeue message
			continue
		}

		// Ensure exactly one row was updated
		if result.RowsAffected != 1 {
			log.Printf("[Worker] Unexpected update count for account (%s), rejecting transaction", transaction.AccountNumber)
			tx.Rollback()
			msg.Reject(false) // Discard message, account might be invalid
			continue
		}
		log.Printf("[Worker] Updated balance for account: %s", transaction.AccountNumber)

		// Commit the transaction to PostgreSQL
		if err := tx.Commit().Error; err != nil {
			log.Println("[Worker] Failed to commit transaction:", err)
			msg.Nack(false, true) // Requeue message
			continue
		}
		log.Println("[Worker] PostgreSQL transaction committed")

		// Insert transaction log into MongoDB
		_, err := mongoDB.Collection("transactions").InsertOne(context.TODO(), transaction)
		if err != nil {
			log.Println("[Worker] Failed to insert transaction into MongoDB:", err)
			msg.Nack(false, true) // Requeue message
			continue
		}
		log.Println("[Worker] Transaction logged in MongoDB")

		log.Println("[Worker] Processed transaction successfully:", transaction)

		// Acknowledge successful processing
		msg.Ack(false)
	}
}

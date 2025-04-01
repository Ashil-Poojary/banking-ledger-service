package worker

import (
	"encoding/json"
	"log"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/streadway/amqp"
)

// PublishTransaction sends a transaction message to RabbitMQ using an existing channel
func PublishTransaction(transaction models.Transaction, ch *amqp.Channel, queueName string) error {

	if err := transaction.Validate(); err != nil {
		log.Println("Transaction validation failed:", err)
		return err
	}

	body, err := json.Marshal(transaction)
	if err != nil {
		log.Println("Failed to marshal transaction:", err)
		return err
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 🔹 Ensure messages persist even after RabbitMQ restarts
			Body:         body,
		},
	)
	if err != nil {
		log.Println("Failed to publish transaction:", err)
		return err
	}

	log.Println("[Worker] Published transaction:", transaction.ID)
	return nil
}

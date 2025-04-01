package worker

import (
	"encoding/json"
	"log"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/streadway/amqp"
)

// RabbitMQPublisher defines the necessary RabbitMQ functions
type RabbitMQPublisher interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// PublishTransaction sends a transaction message to RabbitMQ
func PublishTransaction(transaction models.Transaction, publisher RabbitMQPublisher, queueName string) error {
	if err := transaction.Validate(); err != nil {
		log.Println("Transaction validation failed:", err)
		return err
	}

	body, err := json.Marshal(transaction)
	if err != nil {
		log.Println("Failed to marshal transaction:", err)
		return err
	}

	err = publisher.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("Failed to publish transaction:", err)
		return err
	}

	log.Println("Published transaction:", string(body))
	return nil
}

package worker

import (
	"encoding/json"
	"log"

	"github.com/ashil-poojary/banking-ledger-service/models"
	"github.com/streadway/amqp"
)

// PublishTransaction sends a transaction message to RabbitMQ
// PublishTransaction sends a transaction message to RabbitMQ using an existing channel
func PublishTransaction(transaction models.Transaction, ch *amqp.Channel, queueName string) error {
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

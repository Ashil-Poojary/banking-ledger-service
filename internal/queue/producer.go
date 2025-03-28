package queue

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	_, err = ch.QueueDeclare(
		"transactions_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}
}

func PublishTransaction(msg string) {
	err := ch.Publish(
		"",
		"transactions_queue",
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(msg)},
	)
	if err != nil {
		log.Println("Failed to publish message:", err)
	}
}

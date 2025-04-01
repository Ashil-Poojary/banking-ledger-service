package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var (
	RabbitMQConn    *amqp.Connection
	RabbitMQChannel *amqp.Channel
)

// InitRabbitMQ initializes RabbitMQ connection and a channel
func InitRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	rabbitMQURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a RabbitMQ channel: %v", err)
	}

	log.Println("Connected to RabbitMQ successfully")

	// Store globally
	RabbitMQConn = conn
	RabbitMQChannel = ch
	return conn, ch
}

// CloseRabbitMQ properly closes the connection and channel
func CloseRabbitMQ() {
	if RabbitMQChannel != nil {
		RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
	log.Println("RabbitMQ connection closed")
}

package utils

import (
	"log"

	"github.com/streadway/amqp"
)

// MockRabbitMQ simulates RabbitMQ publishing for tests.
type MockRabbitMQ struct{}

// Publish simulates publishing a message to RabbitMQ.
func (m *MockRabbitMQ) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	log.Println("MockRabbitMQ: Publishing message:", string(msg.Body))
	return nil // Simulate successful publishing
}

package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

	// Declare the queue before consuming
	_, err = ch.QueueDeclare(
		"transactions_queue", // Queue name
		true,                 // Durable
		false,                // Auto-delete
		false,                // Exclusive
		false,                // No-wait
		nil,                  // Arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		"transactions_queue",
		"",
		true,  // Auto-acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Args
	)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	// Process incoming messages
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Received transaction: %s\n", d.Body)
		}
	}()

	fmt.Println("Worker is running... Waiting for messages.")
	<-forever
}

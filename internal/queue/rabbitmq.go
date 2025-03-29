package queue

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

var RabbitMQConn *amqp.Connection

func ProcessTransactionQueue() {
	var err error

	for i := 0; i < 10; i++ {
		RabbitMQConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			log.Println("RabbitMQ connected successfully!")
			break
		}

		log.Printf("Failed to connect to RabbitMQ (attempt %d/10): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after 10 attempts: %v", err)
	}

	ch, err := RabbitMQConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel: ", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("transactions", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare a queue: ", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register a consumer: ", err)
	}

	for msg := range msgs {
		log.Printf("Received a transaction message: %s", msg.Body)
	}
}

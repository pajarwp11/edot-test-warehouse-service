package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp091.Connection

func Connect() {
	var err error
	RabbitConn, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
}

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	rabbitConn *amqp091.Connection
}

func NewRabbitConsumer(rabbitConn *amqp091.Connection) *RabbitConsumer {
	return &RabbitConsumer{
		rabbitConn: rabbitConn,
	}
}

func (r *RabbitConsumer) ConsumeEvents() {
	ch, err := r.rabbitConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exhangeName,
		"topic",
		true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	topics := map[string]string{
		"stock.transfer": "queue_stock_transfer",
	}

	for routingKey, queueName := range topics {
		go startConsumer(ch, queueName, routingKey)
	}

	log.Println("consumer started")
	select {}
}

func startConsumer(ch *amqp091.Channel, queueName, routingKey string) {
	q, err := ch.QueueDeclare(
		queueName,
		true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, routingKey, exhangeName, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {
			var event Event
			json.Unmarshal(d.Body, &event)
			log.Printf("Received from %s: %+v\n", queueName, event)
			handleEvent(event)
		}
	}()
}

func handleEvent(event Event) {
	switch event.Type {
	case "stock.transfer":
		fmt.Println("Handling stock reservation:", event.Data)
	default:
		fmt.Println("Unknown event:", event.Type)
	}
}

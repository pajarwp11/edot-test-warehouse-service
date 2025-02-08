package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

var (
	RabbitConn  *amqp091.Connection
	exhangeName = "stock_events"
)

func Connect() {
	var err error
	RabbitConn, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
}

func ConsumeEvents() {
	ch, err := RabbitConn.Channel()
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

func PublishEvent(eventType string, data map[string]interface{}) error {
	ch, err := RabbitConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exhangeName,
		"topic",
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	event := Event{
		Type: eventType,
		Data: data,
	}
	body, _ := json.Marshal(event)

	return ch.PublishWithContext(
		context.Background(),
		"stock_events",
		eventType,
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

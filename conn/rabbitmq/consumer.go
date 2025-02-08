package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"warehouse-service/entity"

	"github.com/rabbitmq/amqp091-go"
)

type StockHandler interface {
	TransferStock(data interface{}) error
	AddStock(data interface{}) error
	DeductStock(data interface{}) error
}

type RabbitConsumer struct {
	rabbitConn   *amqp091.Connection
	stockHandler StockHandler
}

func NewRabbitConsumer(rabbitConn *amqp091.Connection, stockHandler StockHandler) *RabbitConsumer {
	return &RabbitConsumer{
		rabbitConn:   rabbitConn,
		stockHandler: stockHandler,
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
		go r.startConsumer(ch, queueName, routingKey)
	}

	log.Println("consumer started")
	select {}
}

func (r *RabbitConsumer) startConsumer(ch *amqp091.Channel, queueName, routingKey string) {
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

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {
			var event Event
			json.Unmarshal(d.Body, &event)
			log.Printf("Received from %s: %+v\n", queueName, event)
			err := r.handleEvent(event)
			if err != nil {
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()
}

func (r *RabbitConsumer) handleEvent(event Event) error {
	switch event.Type {
	case entity.StockTransferEvent:
		return r.stockHandler.TransferStock(event.Data)
	case entity.StockAddEvent:
		return r.stockHandler.AddStock(event.Data)
	case entity.StockDeductEvent:
		return r.stockHandler.DeductStock(event.Data)
	default:
		fmt.Println("Unknown event:", event.Type)
		return nil
	}
}

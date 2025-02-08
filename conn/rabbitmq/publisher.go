package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	rabbitConn *amqp091.Connection
}

func NewRabbitPublisher(rabbitConn *amqp091.Connection) *RabbitPublisher {
	return &RabbitPublisher{
		rabbitConn: rabbitConn,
	}
}

func (r *RabbitPublisher) PublishEvent(eventType string, data interface{}) error {
	ch, err := r.rabbitConn.Channel()
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
		exhangeName,
		eventType,
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

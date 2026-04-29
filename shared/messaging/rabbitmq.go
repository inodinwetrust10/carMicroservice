// Package messaging
package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq")
		return nil, fmt.Errorf("failed to connect to rabbitmq: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}
	rmq := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}
	if err := rmq.setupExchangesAndQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed to setup exchanges and queue: %v", err)

	}

	return rmq, nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message string) error {
	return r.Channel.PublishWithContext(ctx,
		"",
		"hello",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	_, err := r.Channel.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.Channel != nil {
		r.Channel.Close()
	}
}

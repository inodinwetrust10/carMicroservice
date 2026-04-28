// Package messaging
package messaging

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp091.Connection
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq")
		return nil, fmt.Errorf("failed to connect to rabbitmq: %v", err)
	}

	return &RabbitMQ{
		conn: conn,
	}, nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

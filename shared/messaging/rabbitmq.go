// Package messaging
package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"ride-sharing/shared/contracts"

	"github.com/rabbitmq/amqp091-go"
)

const (
	TripExchange = "trip"
)

type RabbitMQ struct {
	conn    *amqp091.Connection
	Channel *amqp091.Channel
}
type MessageHandler func(context.Context, amqp091.Delivery) error

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

func (r *RabbitMQ) ConsumeMessages(queueName string, handler MessageHandler) error {
	err := r.Channel.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set Qos: %v", err)
	}
	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)

			if err := handler(ctx, msg); err != nil {
				log.Printf("ERROR: Failed to handle message: %v. Message body: %s", err, msg.Body)
				if nackErr := msg.Nack(false, false); nackErr != nil {
					log.Printf("ERROR: Failed to Nack message: %v", nackErr)
				}
				continue
			}
			if ackErr := msg.Ack(false); ackErr != nil {
				log.Printf("ERROR: Failed to Ack message: %v. Message body: %s", ackErr, msg.Body)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	log.Printf("Publishing message with routing key :%s", routingKey)
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal the message %v", err)
	}
	return r.Channel.PublishWithContext(ctx,
		TripExchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "text/plain",
			Body:         jsonMsg,
			DeliveryMode: amqp091.Persistent,
		})
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	err := r.Channel.ExchangeDeclare(TripExchange, "topic", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %s: %v", TripExchange, err)
	}
	if err := r.declareAndBindQueue(FindAvailableDriversQueue, []string{contracts.TripEventCreated, contracts.TripEventDriverNotInterested}, TripExchange); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) declareAndBindQueue(queueName string, messageTypes []string, exchange string) error {
	q, err := r.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for _, msg := range messageTypes {
		if err := r.Channel.QueueBind(q.Name, msg, exchange, false, nil); err != nil {
			return fmt.Errorf("failed to bind queue to %s: %v", queueName, err)
		}
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

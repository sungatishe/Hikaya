package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"user-list-service/internal/service"

	"github.com/streadway/amqp"
)

// Consumer отвечает за подключение и обработку сообщений из RabbitMQ
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewConsumer создает новое подключение к RabbitMQ
func NewConsumer(url, queueName string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := channel.QueueDeclare(
		queueName, // имя очереди
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &Consumer{conn, channel, queue}, nil
}

// HandleUserEvents подписывается на очередь и обрабатывает события
func (c *Consumer) HandleUserEvents(svc *service.UserListService) error {
	msgs, err := c.channel.Consume(
		c.queue.Name, // имя очереди
		"",           // consumer
		true,         // autoAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Обработка входящих сообщений в отдельной горутине
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var event map[string]interface{}
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			// Определение типа события и вызов соответствующего метода в UserListService
			switch eventType := event["eventType"].(string); eventType {
			case "UserCreated":
				userData := event["data"].(map[string]interface{})
				if err := svc.HandleUserCreated(userData); err != nil {
					log.Printf("Failed to handle 'UserCreated' event: %v", err)
				}
			default:
				log.Printf("Unhandled event type: %s", eventType)
			}
		}
	}()

	log.Printf("Waiting for messages from queue %s", c.queue.Name)
	select {}
}

// Close закрывает соединение и канал RabbitMQ
func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}

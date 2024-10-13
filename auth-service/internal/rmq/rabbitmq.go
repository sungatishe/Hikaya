package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQ создает новое подключение к RabbitMQ
func NewRabbitMQ(url, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Создаем очередь, если она не существует
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

	return &RabbitMQ{conn, channel, queue}, nil
}

// PublishMessage публикует сообщение в очередь
func (r *RabbitMQ) PublishMessage(message string) error {
	err := r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	log.Printf("Message published to queue %s: %s", r.queue.Name, message)
	return nil
}

// Close закрывает подключение и канал RabbitMQ
func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

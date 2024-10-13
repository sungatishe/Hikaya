package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ(url string) *RabbitMQ {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to rmq ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel ", err)
	}

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
	}
}

func (r *RabbitMQ) Publish(queueName string, message []byte) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to declare a queue: ", err)
		return err
	}

	err = r.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Println("Failed to publish a message: ", err)
		return err
	}
	return nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.connection.Close()
}

package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewConsumer(url string) *Consumer {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln("Failed to connect to RMQ: ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open a channel: ", err)
	}

	return &Consumer{
		connection: conn,
		channel:    ch,
	}
}

func (c *Consumer) Consume(queueName string, handleMessage func([]byte)) {
	_, err := c.channel.QueueDeclare(
		queueName, // Имя очереди
		true,      // Durable (выживет после перезапуска RabbitMQ)
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalln("Failed to declare a queue: ", err)
	}

	msgs, err := c.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln("Failed to register a consumer: ", err)
	}

	for msg := range msgs {
		handleMessage(msg.Body)
	}
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.connection.Close()
}

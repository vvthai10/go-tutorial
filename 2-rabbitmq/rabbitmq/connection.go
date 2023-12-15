package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Config struct {
	URL	string
	WaitTime time.Duration
	Attempts int
}

// Connection
type Connection struct {
	ConsumerExchange string
	Config
	Connection *amqp.Connection
	Channel *amqp.Channel
	Delivery <-chan amqp.Delivery
}

// New
func New(consumerExchange string, cfg Config) *Connection {
	conn := &Connection{
		ConsumerExchange: consumerExchange,
		Config: cfg,
	}

	return conn
}

// AttemptConnect
func (c *Connection) AttemptConnect() error {
	var err error
	for i := c.Attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}
		log.Printf("RabbitMQ is trying to connect, attempts left: %d", i)
		time.Sleep(c.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("AttemptConnect - c.connect: %w", err)
	}

	return nil
}

func (c *Connection) connect() error {
	var err error

	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	// Create Channel
	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return fmt.Errorf("c.Connection.Channel: %w", err)
	}

	// Create exchange in channel
	// kind == 'fanout' - find anthor in web
	err = c.Channel.ExchangeDeclare(
		c.ConsumerExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Connection.ExchangeDeclare: %w", err)
	}

	// Create queue
	// Here i set name = "", so when connection, we just can use one queue,
	// I update to create more queue after
	queue, err := c.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Connection.QueueDeclare: %w", err)
	}

	// Add queue to exchange
	err = c.Channel.QueueBind(
		queue.Name,
		"",
		c.ConsumerExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.QueueBind: %w", err)
	}

	// Queue is ready to listen message
	c.Delivery, err = c.Channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.Consume: %w", err)
	}

	return nil
}


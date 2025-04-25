package repository

import (
	"github.com/streadway/amqp"
)

// RabbitMQRepository defines the interface for RabbitMQ operations
type RabbitMQRepository interface {
	ConsumeOrderCreated() (<-chan amqp.Delivery, error)
}

// rabbitMQRepository implements RabbitMQRepository
type rabbitMQRepository struct {
	conn *amqp.Connection
}

// NewRabbitMQRepository creates a new RabbitMQ repository
func NewRabbitMQRepository(conn *amqp.Connection) RabbitMQRepository {
	return &rabbitMQRepository{conn: conn}
}

// ConsumeOrderCreated sets up a consumer for order.created events
func (r *rabbitMQRepository) ConsumeOrderCreated() (<-chan amqp.Delivery, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"order.created", // queue name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return msgs, err
}
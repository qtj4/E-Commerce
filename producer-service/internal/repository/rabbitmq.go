package repository

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"E-Commerce/producer-service/internal/entity"
)

// RabbitMQRepository defines the interface for RabbitMQ operations
type RabbitMQRepository interface {
	PublishOrderCreated(event entity.OrderCreatedEvent) error
}

// rabbitMQRepository implements RabbitMQRepository
type rabbitMQRepository struct {
	conn *amqp.Connection
}

// NewRabbitMQRepository creates a new RabbitMQ repository
func NewRabbitMQRepository(conn *amqp.Connection) RabbitMQRepository {
	return &rabbitMQRepository{conn: conn}
}

// PublishOrderCreated publishes an order.created event to RabbitMQ
func (r *rabbitMQRepository) PublishOrderCreated(event entity.OrderCreatedEvent) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order.created", // queue name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
package repository

import (
	"github.com/streadway/amqp"
)

// RabbitMQRepository defines the interface for RabbitMQ operations
type RabbitMQRepository interface {
    ConsumeOrderCreated() (<-chan amqp.Delivery, error)
    Close() error
}

type rabbitMQRepository struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQRepository(url string) (RabbitMQRepository, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    err = ch.ExchangeDeclare(
        "order_events", // name
        "topic",        // type
        true,          // durable
        false,         // auto-deleted
        false,         // internal
        false,         // no-wait
        nil,          // arguments
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, err
    }

    return &rabbitMQRepository{
        conn:    conn,
        channel: ch,
    }, nil
}

func (r *rabbitMQRepository) ConsumeOrderCreated() (<-chan amqp.Delivery, error) {
    q, err := r.channel.QueueDeclare(
        "order_created_consumer", // name
        true,                    // durable
        false,                   // delete when unused
        false,                   // exclusive
        false,                   // no-wait
        nil,                     // arguments
    )
    if err != nil {
        return nil, err
    }

    err = r.channel.QueueBind(
        q.Name,         // queue name
        "order.created", // routing key
        "order_events", // exchange
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }

    return r.channel.Consume(
        q.Name, // queue
        "",     // consumer
        false,  // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
}

func (r *rabbitMQRepository) Close() error {
    if err := r.channel.Close(); err != nil {
        return err
    }
    return r.conn.Close()
}
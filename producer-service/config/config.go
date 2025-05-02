package config

import (
	"log"

	"github.com/streadway/amqp"
)

type Config struct {
	RabbitMQConn *amqp.Connection
}

func NewConfig() *Config {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	return &Config{
		RabbitMQConn: conn,
	}
}
package config

import (
	"log"

	"github.com/streadway/amqp"
)

// Config holds the configuration for the producer-service
type Config struct {
	RabbitMQConn *amqp.Connection
}

// NewConfig initializes and returns the configuration
func NewConfig() *Config {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	return &Config{
		RabbitMQConn: conn,
	}
}
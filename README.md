# E-Commerce

## Overview
This project is a microservices-based e-commerce system built with Go. It consists of the following services:

- **api-gateway**: Entry point for clients, handles authentication and routes requests to backend services.
- **user-service**: Manages user registration, authentication, and user data.
- **inventory-service**: Manages products and inventory stock.
- **order-service**: Handles order creation and management.
- **producer-service**: Publishes order events to RabbitMQ for asynchronous processing.
- **consumer-service**: Consumes order events from RabbitMQ and updates inventory accordingly.

## How to Run All Services

1. **Start RabbitMQ**  
   Make sure RabbitMQ is running locally (default: `amqp://guest:guest@localhost:5672/`).

2. **Start All Services**  
   You can use the provided script:
   ```sh
   ./start-services.sh
   ```
   Or, start each service in a separate terminal:
   ```sh
   go run ./user-service/cmd/main.go
   go run ./inventory-service/cmd/main.go
   go run ./order-service/cmd/main.go
   go run ./producer-service/cmd/main.go
   go run ./consumer-service/cmd/main.go
   go run ./api-gateway/cmd/main.go
   ```

3. **Ports**
   - user-service: `:50052`
   - inventory-service: `:50055`
   - order-service: `:50051`
   - producer-service: `:50054`
   - api-gateway: (check config, usually `:8080`)

## How to Test Event Flow

1. **Create a User**  
   Use the API gateway to register and log in a user.

2. **Add a Product**  
   Use the API gateway or inventory-service to add products.

3. **Place an Order**  
   Use the API gateway or order-service to create an order.
   - The order-service will notify the producer-service via gRPC.
   - The producer-service publishes an `order.created` event to RabbitMQ.
   - The consumer-service consumes the event and updates inventory via gRPC.

4. **Check Inventory**  
   Query the inventory-service to verify that product stock has been updated.

## Service Descriptions

- **api-gateway**: Handles HTTP requests, authentication, and routes to backend services.
- **user-service**: Manages user accounts and authentication.
- **inventory-service**: Manages products and inventory stock.
- **order-service**: Handles order creation, validation, and notifies producer-service.
- **producer-service**: Receives gRPC notifications from order-service and publishes events to RabbitMQ.
- **consumer-service**: Listens to RabbitMQ for order events and updates inventory accordingly.
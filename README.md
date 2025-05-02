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
   - inventory-service: `:50053`
   - order-service: `:50051`
   - producer-service: `:50054`
   - consumer-service: `50055`
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

---

# API Test Guide

## Authentication & User Tests

### A. Register Admin
- **Method**: `POST`
- **URL**: `http://localhost:8080/auth/register`
- **Headers**:
  - `Content-Type: application/json`
- **Body**:
  ```json
  {
    "email": "admin1234@test.com",
    "password": "admin123"
  }
  ```
- **Check**:
  - Status code should be `201`
  - Copy the token from the response and save it as `admin_token` for later use.

### B. Register User
- **Method**: `POST`
- **URL**: `http://localhost:8080/auth/register`
- **Headers**:
  - `Content-Type: application/json`
- **Body**:
  ```json
  {
    "email": "user@test.com",
    "password": "user123"
  }
  ```
  - Note: If no role is specified, it defaults to "user".
- **Check**:
  - Status code should be `201`
  - Copy the token from the response and save it as `user_token`.

### C. Login as Admin
- **Method**: `POST`
- **URL**: `http://localhost:8080/auth/login`
- **Headers**:
  - `Content-Type: application/json`
- **Body**:
  ```json
  {
    "email": "admin@test.com",
    "password": "admin123"
  }
  ```
- **Check**:
  - Status code should be `200`
  - Copy the token from the response and save it as `admin_token`.

### D. Login as User
- **Method**: `POST`
- **URL**: `http://localhost:8080/auth/login`
- **Headers**:
  - `Content-Type: application/json`
- **Body**:
  ```json
  {
    "email": "user@test.com",
    "password": "user123"
  }
  ```
- **Check**:
  - Status code should be `200`
  - Copy the token from the response and save it as `user_token`.

---

## Product Endpoints (Admin Only)

### A. Create Product
- **Method**: `POST`
- **URL**: `http://localhost:8080/products`
- **Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer {{admin_token}}`
- **Body**:
  ```json
  {
    "name": "Product 1",
    "description": "Description of product 1",
    "price": 10.99,
    "stock": 100,
    "category_id": "category-uuid"
  }
  ```
- **Check**:
  - Status code should be `201`
  - Copy the product ID from the response (e.g., save as `product_id`) for future use.

### B. Update Product
- **Method**: `PATCH`
- **URL**: `http://localhost:8080/products/{{product_id}}`
- **Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer {{admin_token}}`
- **Body**:
  ```json
  {
    "name": "Updated Product 1",
    "price": 12.99
  }
  ```
- **Check**:
  - Status code should be `200`
  - Verify the updated product details in the response.

### C. Delete Product
- **Method**: `DELETE`
- **URL**: `http://localhost:8080/products/{{product_id}}`
- **Headers**:
  - `Authorization: Bearer {{admin_token}}`
- **Check**:
  - Status code should be `200`
  - Confirm deletion by attempting to retrieve the product (should return `404`).

---

## Product Endpoints (User & Admin)

### A. Get Product
- **Method**: `GET`
- **URL**: `http://localhost:8080/products/{{product_id}}`
- **Headers**:
  - `Authorization: Bearer {{user_token}}` (or `admin_token`)
- **Check**:
  - Status code should be `200`
  - Verify the product details in the response.

### B. List Products
- **Method**: `GET`
- **URL**: `http://localhost:8080/products?page=1&page_size=10`
- **Headers**:
  - `Authorization: Bearer {{user_token}}` (or `admin_token`)
- **Check**:
  - Status code should be `200`
  - Verify the list of products and total count in the response.

---

## Order Endpoints (User & Admin)

### A. Create Order (User)
- **Method**: `POST`
- **URL**: `http://localhost:8080/orders`
- **Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer {{user_token}}`
- **Body**:
  ```json
  {
    "items": [
      {
        "product_id": "1ba83505-10b6-49fb-bff6-da5a2ab37554",
        "quantity": 2
      }
    ]
  }
  ```
- **Check**:
  - Status code should be `201`
  - Copy the order ID from the response (e.g., save as `order_id`).

### B. Get Order
- **Method**: `GET`
- **URL**: `http://localhost:8080/orders/{{order_id}}`
- **Headers**:
  - `Authorization: Bearer {{user_token}}` (or `admin_token`)
- **Check**:
  - Status code should be `200`
  - Verify the order details in the response.

### C. List Orders (User)
- **Method**: `GET`
- **URL**: `http://localhost:8080/orders?page=1&page_size=10`
- **Headers**:
  - `Authorization: Bearer {{user_token}}`
- **Check**:
  - Status code should be `200`
  - Verify the list of orders specific to the user.

### D. List Orders (Admin)
- **Method**: `GET`
- **URL**: `http://localhost:8080/orders?user_id={{user_id}}&page=1&page_size=10`
- **Headers**:
  - `Authorization: Bearer {{admin_token}}`
- **Check**:
  - Status code should be `200`
  - Verify the list of orders for the specified `user_id`.

### E. Update Order (Admin)
- **Method**: `PATCH`
- **URL**: `http://localhost:8080/orders/{{order_id}}`
- **Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer {{admin_token}}`
- **Body**:
  ```json
  {
    "status": "completed"
  }
  ```
- **Check**:
  - Status code should be `200`
  - Verify the updated order status in the response.

---

## User Profile Endpoints

### A. Get Current User Profile
- **Method**: `GET`
- **URL**: `http://localhost:8080/users/me`
- **Headers**:
  - `Authorization: Bearer {{user_token}}` (or `admin_token`)
- **Check**:
  - Status code should be `200`
  - Verify the user profile details in the response.

### B. Update Current User Profile (Change Password)
- **Method**: `PUT`
- **URL**: `http://localhost:8080/users/me`
- **Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer {{user_token}}` (or `admin_token`)
- **Body**:
  ```json
  {
    "new_password": "newpassword123"
  }
  ```
- **Check**:
  - Status code should be `200`
  - Verify the updated user profile (e.g., log in with the new password).

---

## Additional Notes
- **Token Usage**: Use the `Authorization: Bearer {{token}}` header for all protected endpoints. Substitute `admin_token` for admin-only actions and `user_token` for user actions.
- **Error Handling**: Watch for `401 Unauthorized` (missing/invalid token) or `403 Forbidden` (insufficient permissions) errors.
- **Prerequisites**:
  - Ensure the PostgreSQL database is running and the `users` table is created.
  - Confirm all services are active on their default ports:
    - `api-gateway`: `:8080`
    - `user-service`: `:50053`
    - `inventory-service`: `:50051`
    - `order-service`: `:50052`
    - `producer-service`: `:50054`
    - `consumer-service`: `50055`

This guide provides a comprehensive manual for testing all endpoints in Postman, ensuring you can validate the full functionality of your API services.
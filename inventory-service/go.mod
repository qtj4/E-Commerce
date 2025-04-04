module github.com/qtj4/E-Commerce/order-service

go 1.21

require (
    github.com/google/uuid v1.3.0
    github.com/jmoiron/sqlx v1.3.5
    github.com/lib/pq v1.10.9
    google.golang.org/grpc v1.58.3
    google.golang.org/protobuf v1.31.0
    github.com/qtj4/E-Commerce/inventory-service/proto v0.0.0
)

replace github.com/qtj4/E-Commerce/inventory-service/proto => ../inventory-service/proto
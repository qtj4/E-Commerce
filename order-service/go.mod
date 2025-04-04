module github.com/yourusername/order-service

go 1.21

require (
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.58.3
	github.com/yourusername/inventory-service/proto v0.0.0
)

replace github.com/yourusername/inventory-service/proto => ../inventory-service/proto
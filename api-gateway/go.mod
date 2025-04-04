module github.com/qtj4/E-Commerce/api-gateway

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/qtj4/E-Commerce/inventory-service/proto v0.0.0
    github.com/qtj4/E-Commerce/order-service/proto v0.0.0
    google.golang.org/grpc v1.58.3
)

replace (
    github.com/qtj4/E-Commerce/inventory-service/proto => ../inventory-service/proto
    github.com/qtj4/E-Commerce/order-service/proto => ../order-service/proto
)
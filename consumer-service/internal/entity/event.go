package entity

// OrderCreatedEvent represents the structure of an order.created event
type OrderCreatedEvent struct {
	OrderID  string   `json:"order_id"`
	Products []string `json:"products"`
}
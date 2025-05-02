package entity

type OrderCreatedEvent struct {
	OrderID  string   `json:"order_id"`
	Products []string `json:"products"`
}
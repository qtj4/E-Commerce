package entity

import (
	"time"

	"github.com/google/uuid"
)

type StockLog struct {
	ID            uuid.UUID `db:"id"`
	ProductID     uuid.UUID `db:"product_id"`
	PreviousStock int       `db:"previous_stock"`
	NewStock      int       `db:"new_stock"`
	ChangeAmount  int       `db:"change_amount"`
	OperationType string    `db:"operation_type"`
	CreatedAt     time.Time `db:"created_at"`
	OrderID       string    `db:"order_id"`
}

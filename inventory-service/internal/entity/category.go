package entity

import (
	"github.com/google/uuid"
)

// Category represents a product category in the inventory system
type Category struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description,omitempty"`
}
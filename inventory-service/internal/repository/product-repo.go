package repository

import (
    "github.com/google/uuid"
    "E-Commerce/inventory-service/internal/entity"
)

type ProductRepository interface {
    Create(p *entity.Product) error
    Get(id uuid.UUID) (*entity.Product, error)
    Update(p *entity.Product) error
    Delete(id uuid.UUID) error
    List(categoryID uuid.UUID, page, pageSize int) ([]*entity.Product, int, error)
    CheckStock(productID uuid.UUID, quantity int) (bool, error)
    UpdateStock(productID uuid.UUID, quantity int) error
}
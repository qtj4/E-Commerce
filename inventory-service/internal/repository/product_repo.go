package repository

import (
	"E-Commerce/inventory-service/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(p *entity.Product) error
	Get(id uuid.UUID) (*entity.Product, error)
	Update(p *entity.Product) error
	Delete(id uuid.UUID) error
	List(categoryID string, page, pageSize int) ([]*entity.Product, int, error)
	CheckStock(productID uuid.UUID, quantity int) (bool, error)
	UpdateStock(productID uuid.UUID, quantity int, orderID string) error
}

type productRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewProductRepository(db *sqlx.DB, redis *redis.Client) ProductRepository {
	return &productRepository{
		db:    db,
		redis: redis,
	}
}

const productCacheKeyPrefix = "product:"
const productCacheDuration = 24 * time.Hour

func (r *productRepository) Create(p *entity.Product) error {
	p.ID = uuid.New()
	_, err := r.db.NamedExec(`
		INSERT INTO products (id, name, description, price, stock, category_id)
		VALUES (:id, :name, :description, :price, :stock, :category_id)`,
		p)
	return err
}

func (r *productRepository) Get(id uuid.UUID) (*entity.Product, error) {
	ctx := context.Background()
	cacheKey := productCacheKeyPrefix + id.String()

	// Try to get from cache first
	cachedProduct, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var product entity.Product
		if err := json.Unmarshal([]byte(cachedProduct), &product); err == nil {
			return &product, nil
		}
	}

	// If not in cache, get from database
	var product entity.Product
	err = r.db.Get(&product, "SELECT * FROM products WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if productJSON, err := json.Marshal(product); err == nil {
		r.redis.Set(ctx, cacheKey, productJSON, productCacheDuration)
	}

	return &product, nil
}

func (r *productRepository) Update(p *entity.Product) error {
	_, err := r.db.NamedExec(`
		UPDATE products 
		SET name = :name, description = :description, price = :price, 
			stock = :stock, category_id = :category_id
		WHERE id = :id`,
		p)

	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	r.redis.Del(ctx, productCacheKeyPrefix+p.ID.String())

	return nil
}

func (r *productRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	r.redis.Del(ctx, productCacheKeyPrefix+id.String())

	return nil
}

func (r *productRepository) List(categoryID string, page, pageSize int) ([]*entity.Product, int, error) {
	var products []*entity.Product
	var total int
	query := "SELECT * FROM products"
	countQuery := "SELECT COUNT(*) FROM products"

	if categoryID != "" {
		query += " WHERE category_id = $1"
		countQuery += " WHERE category_id = $1"
		err := r.db.Get(&total, countQuery, categoryID)
		if err != nil {
			return nil, 0, err
		}
		err = r.db.Select(&products, query+" LIMIT $2 OFFSET $3", categoryID, pageSize, (page-1)*pageSize)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.db.Get(&total, countQuery)
		if err != nil {
			return nil, 0, err
		}
		err = r.db.Select(&products, query+" LIMIT $1 OFFSET $2", pageSize, (page-1)*pageSize)
		if err != nil {
			return nil, 0, err
		}
	}
	return products, total, nil
}

func (r *productRepository) CheckStock(productID uuid.UUID, quantity int) (bool, error) {
	var stock int
	err := r.db.Get(&stock, "SELECT stock FROM products WHERE id = $1", productID)
	if err != nil {
		return false, err
	}
	return stock >= quantity, nil
}

func (r *productRepository) UpdateStock(productID uuid.UUID, quantity int, orderID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Get current stock within transaction
	var currentStock int
	err = tx.Get(&currentStock, "SELECT stock FROM products WHERE id = $1", productID)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %v", err)
	}

	// Calculate new stock
	newStock := currentStock + quantity
	if newStock < 0 {
		return fmt.Errorf("insufficient stock: current=%d, requested=%d", currentStock, -quantity)
	}

	// Update stock
	_, err = tx.Exec("UPDATE products SET stock = $1 WHERE id = $2", newStock, productID)
	if err != nil {
		return fmt.Errorf("failed to update stock: %v", err)
	}

	// Create stock log entry
	stockLog := &entity.StockLog{
		ID:            uuid.New(),
		ProductID:     productID,
		PreviousStock: currentStock,
		NewStock:      newStock,
		ChangeAmount:  quantity,
		OperationType: getOperationType(quantity),
		CreatedAt:     time.Now(),
		OrderID:       orderID,
	}

	_, err = tx.NamedExec(`
		INSERT INTO stock_logs 
		(id, product_id, previous_stock, new_stock, change_amount, operation_type, created_at, order_id)
		VALUES 
		(:id, :product_id, :previous_stock, :new_stock, :change_amount, :operation_type, :created_at, :order_id)`,
		stockLog)
	if err != nil {
		return fmt.Errorf("failed to create stock log: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Invalidate cache after successful update
	ctx := context.Background()
	r.redis.Del(ctx, productCacheKeyPrefix+productID.String())

	log.Printf("Successfully updated stock for product %s: %d → %d (change: %d) [order: %s]",
		productID, currentStock, newStock, quantity, orderID)

	return nil
}

func getOperationType(quantity int) string {
	if quantity > 0 {
		return "RESTOCK"
	}
	return "ORDER"
}

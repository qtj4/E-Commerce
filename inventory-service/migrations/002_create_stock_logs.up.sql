CREATE TABLE stock_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    product_id UUID NOT NULL REFERENCES products (id),
    previous_stock INT NOT NULL,
    new_stock INT NOT NULL,
    change_amount INT NOT NULL,
    operation_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    order_id VARCHAR(255)
);
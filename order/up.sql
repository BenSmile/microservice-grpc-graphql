CREATE TABLE IF NOT EXISTS orders(
    id VARCHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    account_id VARCHAR(27) NOT NULL,
    total_price MONEY NOT NULL
);
CREATE TABLE IF NOT EXISTS order_products(
    order_id VARCHAR(27) REFERENCES orders(id) ON DELETE CASCADE,
    product_id VARCHAR(27) NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (product_id, order_id)
)
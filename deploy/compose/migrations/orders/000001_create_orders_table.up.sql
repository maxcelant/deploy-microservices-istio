CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL, 
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'PENDING'
);

CREATE TABLE order_items (
    order_id INT NOT NULL,
    item_id INT NOT NULL, 
    quantity INT DEFAULT 1 CHECK (quantity > 0),
    PRIMARY KEY (order_id, item_id)
);

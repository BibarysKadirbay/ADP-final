-- order-service/migrations/000001_init.up.sql

CREATE TABLE IF NOT EXISTS orders (
      id UUID PRIMARY KEY,
      user_id TEXT NOT NULL,
      restaurant_id TEXT NOT NULL,
      delivery_id TEXT,
      total_price BIGINT NOT NULL,
      status TEXT NOT NULL,
      payment_status TEXT NOT NULL,
      address TEXT,
      comment TEXT,
      created_at TIMESTAMP DEFAULT NOW()
);
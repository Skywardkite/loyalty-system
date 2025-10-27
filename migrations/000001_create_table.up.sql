CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    login       TEXT NOT NULL UNIQUE,
    password    TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE accounts (
    user_id     INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance     INTEGER NOT NULL DEFAULT 0,
    total_spent INTEGER NOT NULL DEFAULT 0
);
CREATE TYPE order_status AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
);
CREATE TABLE orders (
    id SERIAL           PRIMARY KEY,
    order_number        TEXT NOT NULL UNIQUE,
    user_id             INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status order_status NOT NULL DEFAULT 'NEW',
    points              INTEGER NOT NULL DEFAULT 0,
    uploaded_at         TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
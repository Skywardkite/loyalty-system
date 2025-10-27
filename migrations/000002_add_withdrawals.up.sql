CREATE TABLE withdrawals (
    id SERIAL           PRIMARY KEY,
    order_number        TEXT NOT NULL UNIQUE,
    user_id             INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    points              INTEGER NOT NULL,
    processed_at        TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_withdrawals_user_id ON withdrawals(user_id);
-- Create transactions table (SQLite compatible)
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('purchase', 'generation')),
    amount INTEGER NOT NULL, -- Positive for purchases, negative for generation
    description TEXT,
    external_payment_id TEXT,
    related_template_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indices
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_external_payment_id ON transactions(external_payment_id); 
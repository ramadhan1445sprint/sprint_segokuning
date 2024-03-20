CREATE UNIQUE INDEX IF NOT EXISTS idx_email_phone ON users(email, phone);

CREATE INDEX IF NOT EXISTS idx_created_at_user On users(created_at);
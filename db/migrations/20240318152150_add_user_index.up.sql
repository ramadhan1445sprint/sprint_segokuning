CREATE UNIQUE INDEX IF NOT EXISTS unique_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS unique_phone ON users(phone);
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    comment VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    user_id UUID REFERENCES users(id),
    post_id UUID REFERENCES posts(id)
);

CREATE INDEX IF NOT EXISTS idx_comments ON comments (created_at);
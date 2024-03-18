CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    post_in_html VARCHAR(500) NOT NULL,
    tags VARCHAR[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP,
    user_id UUID REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_tags ON posts USING GIN (tags);
CREATE INDEX IF NOT EXISTS idx_posts ON posts (created_at, post_in_html);
CREATE TABLE comments (
    id UUID PRIMARY KEY,
	comment VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP,
    user_id UUID REFERENCES users(id),
    post_id UUID REFERENCES posts(id)
);

CREATE INDEX idx_comments ON comments (created_at);
CREATE TABLE IF NOT EXISTS friends(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id1 uuid NOT NULL,
  user_id2 uuid NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id1) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id2) REFERENCES users(id) ON DELETE CASCADE
  -- CONSTRAINT check_different_users CHECK (user_id1 <> user_id2)
);

CREATE INDEX IF NOT EXISTS friends_user_id ON friends (user_id1);

CREATE INDEX IF NOT EXISTS current_user_id ON friends (user_id2);
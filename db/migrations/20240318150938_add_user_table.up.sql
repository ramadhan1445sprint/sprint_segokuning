CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE cred_type AS ENUM('phone', 'email');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    credential_type cred_type NOT NULL,
    credential_value VARCHAR(60) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    friend_count int NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
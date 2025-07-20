-- +goose Up
CREATE EXTENSION "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    email TEXT UNIQUE NOT NULL CHECK (email <> ''),
    username VARCHAR(25) UNIQUE NOT NULL CHECK (username <> ''),
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;

DROP EXTENSION "uuid-ossp";
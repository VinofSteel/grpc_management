-- +goose Up
CREATE EXTENSION "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    email TEXT UNIQUE NOT NULL,
    username VARCHAR(25) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;

DROP EXTENSION "uuid-ossp";
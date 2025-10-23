-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
-- +goose Up
CREATE TABLE IF NOT EXISTS Users(
    id UUID PRIMARY KEY,
    created_at timestamp NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

-- +goose Down 
DROP TABLE Users;
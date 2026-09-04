-- +goose Up
ALTER TABLE Users
ADD COLUMN username TEXT NOT NULL UNIQUE;


-- +goose Down
ALTER TABLE Users
DROP COLUMN username;
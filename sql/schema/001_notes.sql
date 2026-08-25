--  +goose Up
CREATE TABLE IF NOT EXISTS Notes(
    note_id UUID PRIMARY KEY,
    daily_note TEXT NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down 
DROP TABLE Notes;
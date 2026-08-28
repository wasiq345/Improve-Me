-- +goose Up
ALTER TABLE Notes
ADD COLUMN user_id UUID NOT NULL 
REFERENCES Users(id) ON DELETE CASCADE;


-- +goose Down
ALTER TABLE Notes
DROP COLUMN user_id;
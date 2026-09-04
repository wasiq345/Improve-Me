-- +goose Up
ALTER TABLE Users
ALTER COLUMN total_notes SET NOT NULL,
ALTER COLUMN current_streak SET NOT NULL,
ALTER COLUMN max_streak SET NOT NULL,
ALTER COLUMN today_count SET NOT NULL;


-- +goose Down
ALTER TABLE Users 
    ALTER COLUMN total_notes DROP NOT NULL, 
    ALTER COLUMN current_streak DROP NOT NULL, 
    ALTER COLUMN max_streak DROP NOT NULL, 
    ALTER COLUMN today_count DROP NOT NULL;

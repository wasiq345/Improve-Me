-- +goose Up
ALTER TABLE Users
ADD COLUMN current_streak SMALLINT default 0,
ADD COLUMN max_streak SMALLINT default 0,
ADD COLUMN total_notes INTEGER default 0,
ADD COLUMN today_count INTEGER default 0;



-- +goose Down
ALTER TABLE Users
DROP COLUMN current_streak,
DROP COLUMN max_streak,
DROP COLUMN total_notes,
DROP COLUMN today_count;
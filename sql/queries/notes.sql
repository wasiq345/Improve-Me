-- name: CreateNote :one
INSERT INTO Notes (note_id, user_id, daily_note, created_at, updated_at) values (
        gen_random_uuid(), $1, $2, NOW(), NOW()
)
RETURNING *;


-- name: DeleteNote :exec
DELETE FROM Notes
where note_id = $1;

-- name: UpdateNote :one
UPDATE Notes
SET daily_note = $2, updated_at = NOW()
where note_id = $1
RETURNING *;

-- name: ReadNote :one
SELECT * FROM Notes WHERE note_id = $1;

-- name: GetAllNotes :many
SELECT * FROM Notes where user_id = $1;

-- name: GetSortedNotes :many
SELECT * FROM Notes WHERE user_id = $1 ORDER BY created_at DESC;
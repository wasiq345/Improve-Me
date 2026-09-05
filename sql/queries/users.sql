-- name: RegisterUser :one
INSERT INTO Users (id, created_at, username, email, password_hash) values (
    gen_random_uuid(), NOW(), $1, $2, $3
)
RETURNING *;


-- name: SearchUserByEmail :one
SELECT * FROM Users where email = $1;

-- name: SearchUserByUserName :one
SELECT * FROM Users where username = $1;

-- name: ResetDailyCount :exec
UPDATE Users SET today_count = 0;

-- name: IncreaseNoteCount :exec
UPDATE Users SET today_count = today_count + 1, total_notes = total_notes + 1 where id = $1;

-- name: UpdateStreak :exec
UPDATE Users SET current_streak = current_streak + 1, max_streak = GREATEST(max_streak, current_streak + 1) WHERE id = $1;


-- name: ResetCurrentStreak :exec
UPDATE Users SET current_streak = 1 WHERE id = $1;
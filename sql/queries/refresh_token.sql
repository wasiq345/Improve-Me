-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, expires_at, revoked_at, userID) values (
    $1, NOW(), NOW(), $2, $3, $4
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

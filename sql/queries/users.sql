-- name: RegisterUser :one
INSERT INTO Users (id, created_at, username, email, password_hash) values (
    gen_random_uuid(), NOW(), $1, $2, $3
)
RETURNING *;


-- name: SearchUserByEmail :one
SELECT * FROM Users where email = $1;

-- name: SearchUserByUserName :one
SELECT * FROM Users where username = $1;
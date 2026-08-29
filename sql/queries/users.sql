-- name: RegisterUser :one
INSERT INTO Users (id, created_at, email, password_hash) values (
    gen_random_uuid(), NOW(), $1, $2
)
RETURNING *;


-- name: SearchUserByEmail :one
SELECT * FROM Users where email = $1;
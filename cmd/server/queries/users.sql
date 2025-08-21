-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id, username, password, created_at;

-- name: GetUserByUsername :one
SELECT id, username, password, created_at
FROM users
WHERE username = $1;

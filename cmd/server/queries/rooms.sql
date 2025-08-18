-- name: CreateRoom :one
INSERT INTO rooms (name, slug)
VALUES ($1, $2)
RETURNING id, name, slug, created_at;

-- name: ListRooms :many
SELECT id, name, slug, created_at
FROM rooms
ORDER BY created_at DESC;

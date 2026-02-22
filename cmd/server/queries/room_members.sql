-- name: JoinRoom :exec
INSERT INTO room_members (room_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: LeaveRoom :exec
DELETE FROM room_members
WHERE room_id = $1 AND user_id = $2;

-- name: ListRoomMembers :many
SELECT u.id, u.username
FROM room_members rm
JOIN users u ON u.id = rm.user_id
WHERE rm.room_id = $1
ORDER BY rm.joined_at;

-- name: GetRoomsForUser :many
SELECT r.id, r.name, r.slug, r.created_at
FROM rooms r
JOIN room_members rm ON rm.room_id = r.id
WHERE rm.user_id = $1
ORDER BY r.created_at DESC;

-- name: IsMember :one
SELECT EXISTS (
  SELECT 1 FROM room_members
  WHERE room_id = $1 AND user_id = $2
) AS is_member;

-- name: CreateMessage :one
INSERT INTO messages (room_id, conversation_id, sender_id, body)
VALUES ($1, $2, $3, $4)
RETURNING id, room_id, conversation_id, sender_id, body, created_at;

-- name: CreateDirectMessage :one
WITH conv AS (
    INSERT INTO conversations (user_a, user_b)
    VALUES (least(@sender_id::bigint, @to_user_id::bigint),
        greatest(@sender_id::bigint, @to_user_id::bigint))
        ON CONFLICT (user_a, user_b)
        DO UPDATE SET user_a = EXCLUDED.user_a
    RETURNING id
)
INSERT INTO messages (conversation_id, sender_id, body)
    SELECT id, @sender_id, @body FROM conv
RETURNING *;

-- name: ListMessagesByRoom :many
SELECT m.id, m.room_id, m.conversation_id, m.sender_id, u.username AS sender_username, m.body, m.created_at
FROM messages m
JOIN users u ON u.id = m.sender_id
WHERE m.room_id = $1
ORDER BY m.created_at DESC
LIMIT $2;

-- name: ListMessagesByConversation :many
SELECT id, room_id, conversation_id, sender_id, body, created_at
FROM messages
WHERE conversation_id = $1
ORDER BY created_at DESC
LIMIT $2;

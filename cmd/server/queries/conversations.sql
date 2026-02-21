-- name: GetOrCreateConversation :one
INSERT INTO conversations (user_a, user_b)
VALUES (LEAST($1, $2), GREATEST($1, $2))
ON CONFLICT (user_a, user_b) DO UPDATE SET user_a = EXCLUDED.user_a
RETURNING id, user_a, user_b;

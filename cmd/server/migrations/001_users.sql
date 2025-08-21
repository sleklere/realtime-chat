-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id            BIGSERIAL PRIMARY KEY,
  username         TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE rooms (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL,
  slug       TEXT UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE room_members (
  room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (room_id, user_id)
);

-- conversaciones 1–1 (user_a < user_b para evitar duplicados)
CREATE TABLE conversations (
  id      BIGSERIAL PRIMARY KEY,
  user_a  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_b  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT conversations_distinct CHECK (user_a <> user_b),
  CONSTRAINT conversations_order CHECK (user_a < user_b),
  UNIQUE (user_a, user_b)
);

-- mensajes (en room o en conversación 1–1)
CREATE TABLE messages (
  id               BIGSERIAL PRIMARY KEY,
  room_id          BIGINT REFERENCES rooms(id) ON DELETE CASCADE,
  conversation_id  BIGINT REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  body             TEXT NOT NULL,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  -- exactamente uno de room_id o conversation_id
  CONSTRAINT messages_target_one CHECK (num_nonnulls(room_id, conversation_id) = 1)
);

-- índices de lectura por historial
CREATE INDEX idx_messages_room_created_at ON messages (room_id, created_at DESC);
CREATE INDEX idx_messages_conv_created_at ON messages (conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender_created_at ON messages (sender_id, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_sender_created_at;
DROP INDEX IF EXISTS idx_messages_conv_created_at;
DROP INDEX IF EXISTS idx_messages_room_created_at;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS room_members;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

// Package ws handles the WebSocket connections for the real-time chat.
package ws

import (
	"encoding/json"
	"time"
)

// WebSocket message types.
const (
	TypeRoomMessage   = "room_message"
	TypeDirectMessage = "direct_message"

	TypeJoinRoom  = "join_room"
	TypeLeaveRoom = "leave_room"

	TypeUserOnline  = "user_online"
	TypeUserOffline = "user_offline"
	TypeUserTyping  = "user_typing"

	TypeLoadRoomHistory  = "load_room_history"
	TypeLoadConversation = "load_conversation"

	TypePing    = "ping"
	TypePong    = "pong"
	TypeError   = "error"
	TypeSuccess = "success"
)

// Message is the envelope for all WebSocket messages.
// Type determines which payload struct to unmarshal into.
type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// RoomMessagePayload is the payload for room messages.
type RoomMessagePayload struct {
	RoomID  int64  `json:"room_id"`
	Content string `json:"content"`
	// fields populated by the server before broadcast
	SenderID       int64  `json:"sender_id,omitempty"`
	SenderUsername string `json:"sender_username,omitempty"`
	MessageID      int64  `json:"message_id,omitempty"`
}

// DirectMessagePayload is the payload for direct (1-to-1) messages.
type DirectMessagePayload struct {
	ToUserID int64  `json:"to_user_id"`
	Content  string `json:"content"`
	// fields populated by the server before broadcast
	FromUserID     int64  `json:"from_user_id,omitempty"`
	FromUsername   string `json:"from_username,omitempty"`
	ConversationID int64  `json:"conversation_id,omitempty"`
	MessageID      int64  `json:"message_id,omitempty"`
}

// JoinRoomPayload is the payload for join/leave room events.
type JoinRoomPayload struct {
	RoomID int64 `json:"room_id"`
}

// UserTypingPayload is the payload for typing indicator events.
type UserTypingPayload struct {
	RoomID   *int64 `json:"room_id,omitempty"`
	ToUserID *int64 `json:"to_user_id,omitempty"`
	IsTyping bool   `json:"is_typing"`
}

// LoadHistoryPayload is the payload for requesting message history.
type LoadHistoryPayload struct {
	RoomID         *int64 `json:"room_id,omitempty"`
	ConversationID *int64 `json:"conversation_id,omitempty"`
	Limit          int    `json:"limit"`
	BeforeID       *int64 `json:"before_id,omitempty"`
}

// ErrorPayload is the payload for error responses sent to the client.
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

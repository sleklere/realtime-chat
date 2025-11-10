package ws

import (
	"encoding/json"
	"time"
)

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

type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type RoomMessagePayload struct {
	RoomID  int64  `json:"room_id"`
	Content string `json:"content"`
	// added in the server
	SenderID       int64  `json:"sender_id,omitempty"`
	SenderUsername string `json:"sender_username,omitempty"`
	MessageID      int64  `json:"message_id,omitempty"`
}

type DirectMessagePayload struct {
	ToUserID int64  `json:"to_user_id"`
	Content  string `json:"content"`
	// added in the server
	FromUserID     int64  `json:"from_user_id,omitempty"`
	FromUsername   string `json:"from_username,omitempty"`
	ConversationID int64  `json:"conversation_id,omitempty"`
	MessageID      int64  `json:"message_id,omitempty"`
}

type JoinRoomPayload struct {
	RoomID int64 `json:"room_id"`
}

type UserTypingPayload struct {
	RoomID   *int64 `json:"room_id,omitempty"`
	ToUserID *int64 `json:"to_user_id,omitempty"`
	IsTyping bool   `json:"is_typing"`
}

type LoadHistoryPayload struct {
	RoomID         *int64 `json:"room_id,omitempty"`
	ConversationID *int64 `json:"conversation_id,omitempty"`
	Limit          int    `json:"limit"`
	BeforeID       *int64 `json:"before_id,omitempty"` // pagination
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

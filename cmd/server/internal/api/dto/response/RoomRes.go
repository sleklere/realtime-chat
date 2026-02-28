package response

import "time"

// RoomRes is the response body for a room.
type RoomRes struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageRes is the response body for a message.
type MessageRes struct {
	ID             int64     `json:"id"`
	RoomID         *int64    `json:"room_id,omitempty"`
	SenderID       int64     `json:"sender_id"`
	SenderUsername string    `json:"sender_username"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
}

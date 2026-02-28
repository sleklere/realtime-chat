package api

import "time"

// AuthRequest represents the login or register request body.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response from login or register endpoints.
type AuthResponse struct {
	User      UserResponse `json:"user"`
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expires_at"`
}

// UserResponse represents a user in API responses.
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// RoomResponse represents a room in API responses.
type RoomResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageResponse represents a message in API responses.
type MessageResponse struct {
	ID             int64     `json:"id"`
	RoomID         *int64    `json:"room_id,omitempty"`
	SenderID       int64     `json:"sender_id"`
	SenderUsername string    `json:"sender_username"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateRoomRequest represents the request body for creating a room.
type CreateRoomRequest struct {
	Name string `json:"name"`
}

// Error represents a structured error response from the server.
type Error struct {
	Code    string `json:"code,omitempty"`
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

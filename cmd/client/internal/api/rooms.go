package api

import "fmt"

// ListRooms returns all available rooms.
func (c *Client) ListRooms() ([]RoomResponse, error) {
	var rooms []RoomResponse
	err := c.do("GET", "/api/v1/rooms", nil, &rooms)
	return rooms, err
}

// CreateRoom creates a new room with the given name.
func (c *Client) CreateRoom(name string) (RoomResponse, error) {
	var room RoomResponse
	err := c.do("POST", "/api/v1/rooms", CreateRoomRequest{Name: name}, &room)
	return room, err
}

// JoinRoom adds the current user to a room.
func (c *Client) JoinRoom(roomID int64) error {
	return c.do("POST", fmt.Sprintf("/api/v1/rooms/%d/join", roomID), nil, nil)
}

// LeaveRoom removes the current user from a room.
func (c *Client) LeaveRoom(roomID int64) error {
	return c.do("DELETE", fmt.Sprintf("/api/v1/rooms/%d/leave", roomID), nil, nil)
}

// GetMessages retrieves messages for a room with the given limit.
func (c *Client) GetMessages(roomID int64, limit int) ([]MessageResponse, error) {
	var messages []MessageResponse
	path := fmt.Sprintf("/api/v1/rooms/%d/messages?limit=%d", roomID, limit)
	err := c.do("GET", path, nil, &messages)
	return messages, err
}

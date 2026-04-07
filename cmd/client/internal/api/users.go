package api

import (
	"fmt"
	"net/url"
)

// GetUserByUsername looks up a user by their username.
func (c *Client) GetUserByUsername(username string) (UserResponse, error) {
	var user UserResponse
	path := fmt.Sprintf("/api/v1/users?username=%s", url.QueryEscape(username))
	if err := c.do("GET", path, nil, &user); err != nil {
		return UserResponse{}, err
	}
	return user, nil
}

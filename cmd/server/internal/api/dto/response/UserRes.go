package response

import "time"

// UserRes represents the response data for a created/logged-in user
type UserRes struct {
	ID        int64 `json:"id"`
	Username  string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

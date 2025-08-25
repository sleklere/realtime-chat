package response

// AuthRes represents the response for a successful registration or login,
// including the issued JWT access token and its expiration time.
type AuthRes struct {
	User      UserRes `json:"user"`
	Token     string  `json:"token"`
	ExpiresAt int64   `json:"expires_at"`
}

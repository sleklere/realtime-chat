package request

// RegisterReq represents the request payload for registering a new user.
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

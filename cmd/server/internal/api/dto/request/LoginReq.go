package request

// LoginReq represents the request payload for authenticating a user
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

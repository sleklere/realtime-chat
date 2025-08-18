package request

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


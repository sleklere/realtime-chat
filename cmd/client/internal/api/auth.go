package api

// Login authenticates with existing credentials and returns a token.
func (c *Client) Login(req AuthRequest) (AuthResponse, error) {
	var res AuthResponse
	err := c.do("POST", "/api/v1/auth/login", req, &res)
	return res, err
}

// Register creates a new account and returns a token.
func (c *Client) Register(req AuthRequest) (AuthResponse, error) {
	var res AuthResponse
	err := c.do("POST", "/api/v1/auth/register", req, &res)
	return res, err
}

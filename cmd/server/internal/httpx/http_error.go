package httpx

import "net/http"

type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code,omitempty"`
	Msg     string `json:"error"`
	Details any    `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *HTTPError) Error() string { return e.Msg }

func New(status int, code, msg string, err error) *HTTPError {
	return &HTTPError{Status: status, Code: code, Msg: msg, Err: err}
}

func BadRequest(code, msg string, err error) *HTTPError {
	return New(http.StatusBadRequest, code, msg, err)
}

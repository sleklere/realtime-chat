package httpx

import (
	"log"
	"net/http"
)

// HTTPError represents a structured error that can be serialized as JSON in HTTP responses.
type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code,omitempty"`
	Msg     string `json:"error"`
	Details any    `json:"details,omitempty"`
	Err     error  `json:"-"`
}

// Error implements the error interface, returning the error message.
func (e *HTTPError) Error() string { return e.Msg }

// New creates a new HTTPError with the given status, code, message, and underlying error.
func New(status int, code, msg string, err error) *HTTPError {
	log.Default().Print("call to http.New")
	return &HTTPError{Status: status, Code: code, Msg: msg, Err: err}
}

// BadRequest creates a new HTTPError with status 400 Bad Request.
func BadRequest(code, msg string, err error) *HTTPError {
	return New(http.StatusBadRequest, code, msg, err)
}

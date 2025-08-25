// Package handlers provides the HTTP endpoint handlers for the API,
// implementing business logic for authentication, system health, and more.
package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	reqdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/request"
	"github.com/sleklere/realtime-chat/cmd/server/internal/auth"
	"github.com/sleklere/realtime-chat/cmd/server/internal/httpx"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authSvc *auth.Service
	logger  *slog.Logger
}

// NewAuthHandler creates a new AuthHandler with the given service and logger.
func NewAuthHandler(s *auth.Service, l *slog.Logger) *AuthHandler {
	return &AuthHandler{authSvc: s, logger: l}
}

// Register handles user registration requests, validates input, and returns a created user or an error.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	h.logger.Debug("register handler")

	var req reqdto.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest("invalid_json", "invalid json", err)
	}

	authRes, err := h.authSvc.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrUsernameTaken) {
			return httpx.New(http.StatusConflict, "username_taken", "username already in use", err)
		}
		return err
	}

	return httpx.JSON(w, http.StatusOK, authRes)
}

// Login handles user login requests (not yet implemented).
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	h.logger.Debug("login handler")

	var req reqdto.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest("invalid_json", "invalid json", err)
	}

	authRes, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		return err
	}

	return httpx.JSON(w, http.StatusOK, authRes)
}

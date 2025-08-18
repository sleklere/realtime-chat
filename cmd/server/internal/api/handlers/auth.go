package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	reqdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/request"
	"github.com/sleklere/realtime-chat/cmd/server/internal/httpx"
	"github.com/sleklere/realtime-chat/cmd/server/internal/user"
)

type AuthHandler struct {
	users *user.Service
	logger *slog.Logger
}

func NewAuthHandler(users *user.Service, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{users: users, logger: logger}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	h.logger.Debug("register handler")

	var req reqdto.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest("invalid_json", "invalid json", err)
	}

	u, err := h.users.Register(context.Background(), req)
	if err != nil {
		if errors.Is(err, user.ErrEmailTaken) {
			return httpx.New(http.StatusConflict, "email_taken", "email already in use", err)
		}
		return err
	}

	return httpx.JSON(w, http.StatusCreated, map[string]any{
		"id":u.ID, "email": u.Email, "created_at": u.CreatedAt,
	})

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("login handler")

}


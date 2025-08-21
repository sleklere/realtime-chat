// Package api contains the HTTP routing layer of the server, including
// middleware, endpoint handlers and helpers
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sleklere/realtime-chat/cmd/server/internal/api/handlers"
	"github.com/sleklere/realtime-chat/cmd/server/internal/auth"
)

// API defines the HTTP API layer with logger and services used by handlers
type API struct {
	Logger      *slog.Logger
	AuthService *auth.Service
}

// RegisterAuthRoutes registers all authentication-related endpoints under /auth
func (a *API) RegisterAuthRoutes(r chi.Router) {
	h := handlers.NewAuthHandler(a.AuthService, a.Logger)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", a.Handle(h.Register))
		r.Post("/login", a.Handle(h.Login))
	})
}

// RegisterSystemRoutes registers system-level endpoints such as health checks
func (a *API) RegisterSystemRoutes(r chi.Router) {
	h := handlers.NewSystemHandler(a.Logger)
	r.Get("/healthz", a.Handle(h.Health))
}

// RespondJSON encodes a value to JSON and writes it to the response with the given status
func (a *API) RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// RespondError sends an error response in JSON format with the given status and message
func (a *API) RespondError(w http.ResponseWriter, status int, msg string) {
	a.RespondJSON(w, status, map[string]string{"error": msg})
}

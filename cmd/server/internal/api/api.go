// Package api contains the HTTP routing layer of the server, including
// middleware, endpoint handlers and helpers
package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/sleklere/realtime-chat/cmd/server/internal/api/handlers"
	"github.com/sleklere/realtime-chat/cmd/server/internal/auth"
)

// API defines the HTTP API layer with logger and services used by handlers
type API struct {
	Logger      *slog.Logger
	AuthService *auth.Service
	AuthConfig  *auth.Config
}

// RegisterAuthRoutes registers all authentication-related endpoints under /auth
func (a *API) registerAuthRoutes(r chi.Router) {
	h := handlers.NewAuthHandler(a.AuthService, a.Logger)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", a.handle(h.Register))
		r.Post("/login", a.handle(h.Login))
	})
}

// RegisterSystemRoutes registers system-level endpoints such as health checks
func (a *API) registerSystemRoutes(r chi.Router) {
	h := handlers.NewSystemHandler(a.Logger)
	r.Get("/healthz", a.handle(h.Health))
}

// Package api contains the HTTP routing layer of the server, including
// middleware, endpoint handlers and helpers
package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/sleklere/realtime-chat/cmd/server/internal/api/handlers"
	"github.com/sleklere/realtime-chat/cmd/server/internal/auth"
	dbstore "github.com/sleklere/realtime-chat/cmd/server/internal/store"
)

// API defines the HTTP API layer with logger and services used by handlers
type API struct {
	Logger      *slog.Logger
	AuthService *auth.Service
	AuthConfig  *auth.Config
	Queries     *dbstore.Queries
}

// RegisterAuthRoutes registers all authentication-related endpoints under /auth
func (a *API) registerAuthRoutes(r chi.Router) {
	h := handlers.NewAuthHandler(a.AuthService, a.Logger)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", a.handle(h.Register))
		r.Post("/login", a.handle(h.Login))
	})
}

// registerRoomRoutes registers all room-related endpoints under /rooms
func (a *API) registerRoomRoutes(r chi.Router) {
	h := handlers.NewRoomHandler(a.Queries, a.Logger)
	r.Route("/rooms", func(r chi.Router) {
		r.Post("/", a.handle(h.Create))
		r.Get("/", a.handle(h.List))
		r.Get("/{slug}", a.handle(h.GetBySlug))
		r.Post("/{roomID}/join", a.handle(h.Join))
		r.Delete("/{roomID}/leave", a.handle(h.Leave))
		r.Get("/{roomID}/messages", a.handle(h.Messages))
	})
}

// registerSystemRoutes registers system-level endpoints such as health checks
func (a *API) registerSystemRoutes(r chi.Router) {
	h := handlers.NewSystemHandler(a.Logger)
	r.Get("/healthz", a.handle(h.Health))
}

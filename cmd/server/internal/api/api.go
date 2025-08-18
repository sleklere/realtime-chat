// Package api contains the HTTP routing layer of the server, including
// middleware, endpoint handlers and helpers
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sleklere/realtime-chat/cmd/server/internal/api/handlers"
	"github.com/sleklere/realtime-chat/cmd/server/internal/user"
)

type API struct {
	Logger *slog.Logger
	Users *user.Service
}

func (a *API) RegisterAuthRoutes(r chi.Router) {
	h := handlers.NewAuthHandler(a.Users, a.Logger)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", a.Handle(h.Register))
		// r.Post("/login", a.Handle(h.Login))
	})
}

func (a *API) RegisterSystemRoutes(r chi.Router) {
	h := handlers.NewSystemHandler(a.Logger)
	r.Get("/healthz", a.Handle(h.Health))
}

func (a *API) RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (a *API) RespondError(w http.ResponseWriter, status int, msg string) {
	a.RespondJSON(w, status, map[string]string{"error": msg})
}

// Package api contains the HTTP routing layer of the server, including
// middleware and endpoint handlers
package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter initializes the HTTP router with default middlewares and base routes.
func NewRouter(a *API) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60e9)) // 1min

	a.RegisterSystemRoutes(r)

	r.Route("/api/v1", func(api chi.Router) {
		a.RegisterAuthRoutes(api)
	})

	return r
}

// command server is the entrypoint for the chat server
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sleklere/realtime-chat/cmd/server/internal/api"
	"github.com/sleklere/realtime-chat/cmd/server/internal/auth"
	"github.com/sleklere/realtime-chat/cmd/server/internal/db"
	dbstore "github.com/sleklere/realtime-chat/cmd/server/internal/store"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	pool, err := db.NewPool(ctx)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("err pinging pool: %v", err)
	}

	authCfg := &auth.Config{
		JWTSecret: []byte(getenv("JWT_SECRET", "default_secret")),
		Issuer:    "realtime-chat",
		AccessTTL: 15 * time.Minute,
	}
	queries := dbstore.New(pool)
	authSvc := auth.NewService(queries, logger, authCfg)

	a := &api.API{
		Logger:      logger,
		AuthService: authSvc,
		AuthConfig:  authCfg,
	}

	addr := ":" + getenv("PORT", "8080")

	r := api.NewRouter(a)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

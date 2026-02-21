// Package config provides configuration loading for the chat client.
package config

import "os"

// Config holds the client configuration values.
type Config struct {
	ServerURL string
	WSURL     string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	return Config{
		ServerURL: getenv("SERVER_URL", "http://localhost:8080"),
		WSURL:     getenv("WS_URL", "ws://localhost:8080"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

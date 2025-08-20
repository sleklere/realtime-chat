package handlers

import (
	"log/slog"
	"net/http"
)

// SystemHandler handles system-related HTTP requests such as health checks.
type SystemHandler struct {
	logger *slog.Logger
}

// NewSystemHandler creates a new SystemHandler with the given logger.
func NewSystemHandler(logger *slog.Logger) *SystemHandler {
	return &SystemHandler{logger: logger}
}

// Health responds with a simple 200 OK and "ok" body to indicate the service is healthy.
func (h *SystemHandler) Health(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		h.logger.Error("write health response", "err", err)
		return err
	}
	return nil
}

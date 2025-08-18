package handlers

import (
	"log/slog"
	"net/http"
)

type SystemHandler struct {
	logger *slog.Logger
}

func NewSystemHandler(logger *slog.Logger) *SystemHandler {
	return &SystemHandler{logger: logger}
}

func (h *SystemHandler) Health(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		h.logger.Error("write health response", "err", err)
		return err
	}
	return nil
}


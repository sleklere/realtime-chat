package httpx

import (
	"encoding/json"
	"net/http"
)

// JSON writes the given value as a JSON response with the specified HTTP status.
// If the status is 204 No Content, it only writes the header without a body.
func JSON(w http.ResponseWriter, status int, v any) error {
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

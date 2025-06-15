package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// UuidHandler returns an HTTP handler that generates and returns a new UUID.
//
// @Summary      Generate UUID
// @Description  Returns a new UUID if the request is within rate limits.
// @Tags         uuid
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "UUID generated"
// @Failure      429  {object}  map[string]string  "Rate limit exceeded"
// @Router       /generate [get]
func UuidHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"uuid": id.String()})
	})
}

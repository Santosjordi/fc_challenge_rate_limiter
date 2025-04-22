package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// uuidHandler returns a handler that can be wrapped with middleware
func UuidHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"uuid": id.String()})
	})
}

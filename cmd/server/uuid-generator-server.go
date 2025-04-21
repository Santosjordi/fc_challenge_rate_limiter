package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// This is a simple UUID generator server that generates UUIDs and returns them in JSON format.
func main() {
	// Start the server
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		// Generate a new UUID
		uuid := uuid.New()

		// Return the UUID in JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"uuid": uuid.String()})
	})

	log.Println("Starting UUID generator server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

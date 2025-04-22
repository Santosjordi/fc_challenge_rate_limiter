package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/santosjordi/fc_challenge_rate_limiter/config"
	"github.com/santosjordi/fc_challenge_rate_limiter/internal/webserver/handlers"
)

func main() {

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config: %+v\n", cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(middleware.RateLimiter)

	r.Get("/generate", handlers.UuidHandler().ServeHTTP)

	log.Println("Starting UUID generator server on :8080")
	// change the mux from the standard library to chi
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

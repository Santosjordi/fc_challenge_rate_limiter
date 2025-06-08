package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/santosjordi/fc_challenge_rate_limiter/config"
	db "github.com/santosjordi/fc_challenge_rate_limiter/internal/infra/ratelimiter"
	"github.com/santosjordi/fc_challenge_rate_limiter/internal/webserver/handlers"
	mw "github.com/santosjordi/fc_challenge_rate_limiter/internal/webserver/middleware"
)

// main starts the UUID generator HTTP server with rate limiting and logging middleware.
// It loads configuration, connects to Redis, sets up the chi router, and listens on port 8080.
func main() {
	log.Println("Initializing context with timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("Loading configuration from .env file...")
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Loaded config: %+v\n", cfg)

	log.Printf("Connecting to Redis at %s:%s...", cfg.RedisHost, cfg.RedisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB, // use default DB
	})

	log.Println("Pinging Redis to verify connection...")
	redisStorage := db.NewRedisStorage(redisClient)
	if err := redisStorage.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis.")

	log.Println("Initializing rate limiter middleware...")
	rateLimiter := mw.NewRateLimitMiddleware(redisStorage, cfg)

	log.Println("Setting up HTTP router and middleware...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
	r.Use(rateLimiter.Handler)

	log.Println("Registering /generate endpoint handler...")
	r.Get("/generate", handlers.UuidHandler().ServeHTTP)

	log.Printf("Starting UUID generator server on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

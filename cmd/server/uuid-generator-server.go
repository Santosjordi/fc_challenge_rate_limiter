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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config: %+v\n", cfg)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB, // use default DB
	})

	redisStorage := db.NewRedisStorage(redisClient)
	if err := redisStorage.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	rateLimiter := mw.NewRateLimitMiddleware(redisStorage, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(rateLimiter.Handler)

	r.Get("/generate", handlers.UuidHandler().ServeHTTP)

	log.Println("Starting UUID generator server on :8080")
	// change the mux from the standard library to chi
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

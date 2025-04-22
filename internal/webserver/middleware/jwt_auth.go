package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/santosjordi/fc_challenge_rate_limiter/config"
)

// JWTAuthMiddleware will check the API_KEY header for a valid JWT token
func JWTAuthMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the API_KEY header
			apiKey := r.Header.Get("API_KEY")
			if apiKey == "" {
				http.Error(w, "API_KEY header is missing", http.StatusUnauthorized)
				return
			}

			// Parse the JWT token
			_, err := jwtauth.VerifyToken(cfg.TokenAuth, apiKey)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
				return
			}

			// Proceed with the next handler if token is valid
			next.ServeHTTP(w, r)
		})
	}
}

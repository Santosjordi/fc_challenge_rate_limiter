package middleware

import (
	"net/http"

	"github.com/santosjordi/fc_challenge_rate_limiter/internal/infra/db"
	"github.com/santosjordi/fc_challenge_rate_limiter/internal/webserver/utils"
)

func RateLimit(next http.Handler) http.Handler {
	// TODO: The rate limiter should block requests from the same IP address if
	// the rate at which the requests are made exceeds a threshold established in the .env file.
	// Q: How does the rate limiter tracks the IPs? A: Through querying redis for a lockout key.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		storage := db.RateLimiterStorage

		// Determine key (token or IP)
		key := utils.GetRateLimitKey(r)

		// Check if this key is locked out
		locked, err := storage.IsLockedOut(ctx, key)
		if err != nil {
			// Fail open: allow the request if Redis has issues
			http.Error(w, "Internal rate limit check error", http.StatusInternalServerError)
			return
		}

		if locked {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		// Record the request
		err = storage.RegisterRequest(ctx, key)
		if err != nil {
			// Optional: log error or respond accordingly
			http.Error(w, "Error registering request", http.StatusInternalServerError)
			return
		}

		// Proceed to next handler
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/santosjordi/fc_challenge_rate_limiter/config"
	db "github.com/santosjordi/fc_challenge_rate_limiter/internal/infra/ratelimiter"
	"github.com/santosjordi/fc_challenge_rate_limiter/internal/webserver/utils"
)

type RateLimitMiddleware struct {
	storage db.RateLimiter
	config  *config.Config
}

func NewRateLimitMiddleware(storage db.RateLimiter, config *config.Config) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		storage: storage,
		config:  config,
	}
}

// Handler returns an HTTP middleware that enforces rate limiting based on the request's key.
// It checks if the request is allowed using the configured storage backend. If the rate limit
// is exceeded, it sets a lockout period and responds with HTTP 429 (Too Many Requests), including
// rate limit headers such as X-RateLimit-Reset. If allowed, it adds X-RateLimit-Remaining and
// X-RateLimit-Limit headers to the response. The rate limits and lockout durations are determined
// based on whether the key is IP-based or token-based.
func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		key := utils.GetRateLimitKey(r)

		allowed, remaining, err := m.storage.CheckAndIncrement(ctx, key)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "rate limit check failed")
			return
		}

		if !allowed {
			// Use appropriate lockout duration based on key type
			lockoutDuration := m.getLockoutDuration(key)
			m.storage.SetLockOut(ctx, key, lockoutDuration)

			w.Header().Set("X-RateLimit-Reset", time.Now().Add(lockoutDuration).Format(time.RFC3339))
			utils.RespondWithError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		// Add rate limit headers with appropriate limits based on key type
		maxRequests := m.config.IPLimitPerSecond
		if utils.IsTokenBasedKey(key) {
			maxRequests = m.config.TokenLimitPerSecond
		}

		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(maxRequests))

		next.ServeHTTP(w, r)
	})
}

func (m *RateLimitMiddleware) getLockoutDuration(key string) time.Duration {
	if utils.IsTokenBasedKey(key) {
		return m.config.TokenLockoutDuration
	}
	return m.config.IPLockoutDuration
}

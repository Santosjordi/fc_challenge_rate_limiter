package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetRateLimitKey(r *http.Request) string {
	clientIP := GetClientIP(r)
	apiKey := r.Header.Get("API_KEY")

	// API_KEY supercedes the client IP address if the API_KEY is not empty
	if apiKey == "" {
		return clientIP
	}
	return apiKey
}

func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For may contain a comma-separated list of IPs
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0]) // the client IP is the first in the list
	}

	// Fallback to X-Real-IP
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // return as-is if splitting fails
	}
	return host
}

// Check this back later for a less lazy implementation
// type TokenInfo struct {
// 	Valid  bool
// 	Claims map[string]interface{}
// }

// func ParseOptionalToken(r *http.Request, cfg *config.Config) TokenInfo {
// 	tokenStr := r.Header.Get("API_KEY")
// 	if tokenStr == "" {
// 		return TokenInfo{Valid: false}
// 	}

// 	_, err := jwtauth.VerifyToken(cfg.TokenAuth, tokenStr)
// 	if err != nil {
// 		return TokenInfo{Valid: false}
// 	}

// 	// // If claims are in correct format, we return them
// 	// if mapClaims, ok := claims.(map[string]interface{}); ok {
// 	// 	return TokenInfo{
// 	// 		Valid:  true,
// 	// 		Claims: mapClaims,
// 	// 	}
// 	// }

// 	return TokenInfo{Valid: false}
// }

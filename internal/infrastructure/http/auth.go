package http

import (
	"crypto/subtle"
	"net/http"


)

// APIKeyHeader is the HTTP header clients must carry the shared secret in.
const APIKeyHeader = "X-API-Key"

// AuthMiddleware guards mutating routes with an API-key check. It requires the
// X-API-Key header to be present and to match the configured key using a
// constant-time comparison. When no key is configured the server fails closed:
// every request is rejected until a key is set.
//
// The key is read once at build time into apiKey so the middleware is a tight,
// allocation-free hot path.
func AuthMiddleware(apiKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			provided := r.Header.Get(APIKeyHeader)

			// Fail closed: no configured key means no mutations are allowed.
			if apiKey == "" || provided == "" {
				writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
				return
			}

			if subtle.ConstantTimeCompare([]byte(provided), []byte(apiKey)) != 1 {
				writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package http

import (
	"net/http"
	"sync"
	"time"
)

// throttleWindow is the fixed window over which per-client-IP mutation requests
// are counted.
const throttleWindow = time.Minute

// rateLimiter is a fixed-window, per-client-IP counter. It rejects requests
// exceeding the configured limit per window with 429 before the request reaches
// the handler, defending the API key against brute force and throttling
// upload/delete abuse.
type rateLimiter struct {
	limit  int
	window time.Duration

	mu      sync.Mutex
	buckets map[string]*bucket
}

// bucket holds the per-IP counter state for the current window.
type bucket struct {
	windowStart time.Time
	count       int
}

// newRateLimiter builds a fixed-window per-IP limiter. A limit <= 0 returns nil
// so callers skip registering the middleware and throttling stays disabled.
func newRateLimiter(limit int) *rateLimiter {
	if limit <= 0 {
		return nil
	}
	return &rateLimiter{
		limit:   limit,
		window:  throttleWindow,
		buckets: make(map[string]*bucket),
	}
}

// Middleware returns a chi-compatible middleware that rejects requests exceeding
// the per-IP limit with 429 before invoking next.
func (rl *rateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.allow(r.RemoteAddr) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// allow records a request for ip and reports whether it may proceed: the first
// `limit` requests in the current window proceed; later ones do not.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b := rl.buckets[ip]
	now := time.Now()
	if b == nil || now.Sub(b.windowStart) >= rl.window {
		rl.buckets[ip] = &bucket{windowStart: now, count: 1}
		return true
	}
	b.count++
	return b.count <= rl.limit
}

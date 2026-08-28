package http

import (
	"net/http"


)

// cacheControl returns a middleware that sets the Cache-Control response
// header for every request that reaches the wrapped handlers.
func cacheControl(header string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", header)
			next.ServeHTTP(w, r)
		})
	}
}

// pageCache is the caching policy for HTML documents: always revalidate,
// so updated pages are never served stale.
func pageCache() func(next http.Handler) http.Handler {
	return cacheControl("no-cache")
}

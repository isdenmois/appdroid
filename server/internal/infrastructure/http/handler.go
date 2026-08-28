// Package http implements the delivery layer: HTTP routes and handlers built
// on the chi v5 router.
package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// throttleDefault is the per-client-IP ceiling on mutating API requests when
// h.throttleLimit is unset (production). It is enforced by the fixed-window
// per-IP rate limiter in ratelimit.go.
const throttleDefault = 10

// writeJSON marshals v and writes it as a JSON response with the given status.
//
// json.Marshal + w.Write is used deliberately instead of json.Encoder.Encode:
// the encoder appends a trailing newline, which would break response bodies
// that must equal exactly e.g. "pong" or "[]".
func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

// Handler bundles the delivery-layer dependencies.
type Handler struct {
	// Apps is the application-layer service the handlers delegate to.
	Apps *AppsHandler
	// Files serves stored APK files.
	Files *FilesHandler
	// Pages renders the SSR pages.
	Pages *PagesHandler
	// APIKey is the shared secret guarding mutating routes. Empty => fail
	// closed.
	APIKey string
	// throttleLimit caps mutating requests per client IP via the fixed-window
	// rate limiter. Zero => New applies the production default of 10/min;
	// tests inject a small value for deterministic rate-limit assertions.
	throttleLimit int
}

// New creates the chi router and registers all routes.
//
// The static admin frontend is embedded in the binary and served by the
// catch-all handler; see static.go.
func New(h *Handler) *chi.Mux {
	// Production throttle limit, overridable by tests via h.throttleLimit.
	limit := throttleDefault
	if h.throttleLimit > 0 {
		limit = h.throttleLimit
	}

	r := chi.NewRouter()
 	r.Use(middleware.Logger, middleware.Recoverer)

	// Static admin frontend: embedded files served by the catch-all at their
	// URL path; "/" serves index.html. The catch-all is wrapped so SSR pages
	// registered below always take precedence.
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		serveStatic(w, req)
	})

	// APK files.
	r.Get("/file/{file}", h.Files.Get)
	r.Get("/apk/{file}/{hash}", h.Files.Get)

	// Server-rendered pages (Obtainium): always revalidate, never cached.
	r.Group(func(r chi.Router) {
		r.Use(pageCache())
		r.Get("/apps", h.Pages.AppList)
		r.Get("/app/{id}", h.Pages.AppPage)
	})

	// JSON API.
	r.Route("/api", func(r chi.Router) {
		// Read-only endpoints: health check and the Obtainium-compatible list.
		r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
			io.WriteString(w, "pong")
		})
		r.Get("/list", h.Apps.List)

		// Mutating routes are guarded by the API-key middleware and
		// rate-limited by client IP to defend the shared secret against
		// brute force and throttle upload/delete abuse.
		r.Group(func(r chi.Router) {
			if rl := newRateLimiter(limit); rl != nil {
				r.Use(rl.Middleware())
			}
			r.Use(AuthMiddleware(h.APIKey))
			r.Post("/upload", h.Apps.Upload)
			r.Delete("/{id}", h.Apps.Delete)
		})
	})

	return r
}


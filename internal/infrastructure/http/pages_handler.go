package http

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	appsvc "github.com/isdenmois/appdroid/internal/application/app"
	domainapp "github.com/isdenmois/appdroid/internal/domain/app"
	"github.com/isdenmois/appdroid/internal/infrastructure/http/views"
)

// PagesHandler renders the server-side HTML pages.
type PagesHandler struct {
	svc      *appsvc.Service
	renderer *views.Renderer
}

// NewPagesHandler creates a PagesHandler.
func NewPagesHandler(svc *appsvc.Service, renderer *views.Renderer) *PagesHandler {
	return &PagesHandler{svc: svc, renderer: renderer}
}

// AppList renders the /apps page listing all apps as Obtainium links.
func (h *PagesHandler) AppList(w http.ResponseWriter, r *http.Request) {
	apps, err := h.svc.List(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	baseURL := baseURLFromHost(r.Host)
	body, err := h.renderer.AppList(baseURL, apps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(body)
}

// AppPage renders the /app/:id page with a download link for the app.
func (h *PagesHandler) AppPage(w http.ResponseWriter, r *http.Request) {
	app, err := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := h.renderer.AppPage([]domainapp.App{*app})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(body)
}

// baseURLFromHost derives the base URL the same way the previous
// implementation did: plain http when the host starts with "192" (local
// network), https otherwise.
func baseURLFromHost(host string) string {
	scheme := "https"
	if strings.HasPrefix(host, "192") {
		scheme = "http"
	}
	return scheme + "://" + host
}

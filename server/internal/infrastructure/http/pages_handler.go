package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appsvc "github.com/isdenmois/appdroid/server/internal/application/app"
	domainapp "github.com/isdenmois/appdroid/server/internal/domain/app"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/http/views"
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
func (h *PagesHandler) AppList(c *gin.Context) {
	apps, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	baseURL := baseURLFromHost(c.Request.Host)
	body, err := h.renderer.AppList(baseURL, apps)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", body)
}

// AppPage renders the /app/:id page with a download link for the app.
func (h *PagesHandler) AppPage(c *gin.Context) {
	app, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	body, err := h.renderer.AppPage([]domainapp.App{*app})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", body)
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

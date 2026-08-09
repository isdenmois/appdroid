// Package http implements the delivery layer: Gin routes and handlers.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler bundles the delivery-layer dependencies.
type Handler struct {
	// Router is the configured Gin engine.
	Router *gin.Engine
	// Apps is the application-layer service the handlers delegate to.
	Apps *AppsHandler
	// Files serves stored APK files.
	Files *FilesHandler
	// Pages renders the SSR pages.
	Pages *PagesHandler
}

// New creates the Gin engine and registers all routes.
//
// The static admin frontend is embedded in the binary and served by the
// catch-all handler; see static.go.
func New(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Static admin frontend: every embedded file is served by the catch-all
	// at its URL path; "/" serves index.html. The NoRoute handler is wrapped
	// in pageCache so the SSR pages registered below always take precedence.
	r.NoRoute(pageCache(), serveStatic)

	// APK files.
	r.GET("/file/:file", h.Files.Get)
	r.GET("/apk/:file/:hash", h.Files.Get)

	// Server-rendered pages (Obtainium). Always revalidate, never cached.
	r.Group("/", pageCache()).GET("/apps", h.Pages.AppList)
	r.Group("/", pageCache()).GET("/app/:id", h.Pages.AppPage)

	// JSON API.
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
		api.GET("/list", h.Apps.List)
		api.POST("/upload", h.Apps.Upload)
		api.DELETE("/:id", h.Apps.Delete)
	}

	return r
}

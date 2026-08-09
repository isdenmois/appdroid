package http

import (
	"embed"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFS embed.FS

// staticType maps the extensions of the embedded assets to their media types.
// A fixed table keeps the served Content-Type independent of the host system
// (the production image is scratch and has no /etc/mime.types).
var staticType = map[string]string{
	".html": "text/html; charset=utf-8",
	".js":   "text/javascript; charset=utf-8",
	".css":  "text/css; charset=utf-8",
}

// serveStatic is the catch-all handler serving the embedded admin frontend:
// "/" resolves to index.html and any other path resolves to the embedded file
// of the same name. Unknown paths are a plain 404.
func serveStatic(c *gin.Context) {
	// Only the document methods are served; anything else (POST, PUT, ...)
	// that reaches the catch-all is a 404 just like before.
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		c.Status(http.StatusNotFound)
		return
	}

	path := strings.TrimPrefix(c.Request.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	data, err := staticFS.ReadFile("static/" + path)
	if err != nil {
		c.Status(http.StatusNotFound) // missing file or a directory
		return
	}

	// HTML entry documents always revalidate; JS/CSS subresources may be
	// cached for a day in release mode. In dev mode (gin debug) nothing is
	// cached, so a local rebuild is always picked up. No validators are
	// emitted: embedded files all report the Unix epoch as their modification
	// time, so Last-Modified/ETag would be wrong (content changes across
	// rebuilds, the timestamp does not).
	if gin.Mode() == gin.DebugMode {
		c.Header("Cache-Control", "no-cache")
	} else if filepath.Ext(path) == ".html" {
		c.Header("Cache-Control", "no-cache")
	} else {
		c.Header("Cache-Control", "public, max-age=86400")
	}

	typ := staticType[filepath.Ext(path)]
	if typ == "" {
		typ = mime.TypeByExtension(filepath.Ext(path))
		if typ == "" {
			typ = "application/octet-stream"
		}
	}
	c.Data(http.StatusOK, typ, data)
}

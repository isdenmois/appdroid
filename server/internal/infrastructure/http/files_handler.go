package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appsvc "github.com/isdenmois/appdroid/server/internal/application/app"
)

// FilesHandler serves the stored APK files.
type FilesHandler struct {
	svc *appsvc.Service
}

// NewFilesHandler creates a FilesHandler.
func NewFilesHandler(svc *appsvc.Service) *FilesHandler {
	return &FilesHandler{svc: svc}
}

// Get serves the file named by the :file path param from the data directory.
func (h *FilesHandler) Get(c *gin.Context) {
	name := c.Param("file")
	f, err := h.svc.OpenFile(c.Request.Context(), name)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	defer f.Close()

	c.DataFromReader(http.StatusOK, -1, "application/vnd.android.package-archive", f, nil)
}

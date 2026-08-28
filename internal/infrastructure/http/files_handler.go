package http

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	appsvc "github.com/isdenmois/appdroid/internal/application/app"
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
func (h *FilesHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "file")
	f, err := h.svc.OpenFile(r.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()

	// No Content-Length / no range handling: matches the old
	// DataFromReader(...,-1,...) chunked behavior.
	w.Header().Set("Content-Type", "application/vnd.android.package-archive")
	io.Copy(w, f)
}

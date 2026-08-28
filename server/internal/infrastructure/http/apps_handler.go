package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	appsvc "github.com/isdenmois/appdroid/server/internal/application/app"
	domainapp "github.com/isdenmois/appdroid/server/internal/domain/app"
)

// appDTO is the JSON representation of an app exposed by the API. The field
// names match the previous implementation exactly.
type appDTO struct {
	ID          string `json:"id"`
	AppID       string `json:"appId"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	VersionName string `json:"versionName"`
	Type        string `json:"type"`
	Apk         string `json:"apk"`
}

func toDTO(a domainapp.App) appDTO {
	return appDTO{
		ID:          a.ID,
		AppID:       a.AppID,
		Name:        a.Name,
		Version:     a.Version,
		VersionName: a.VersionName,
		Type:        string(a.Type),
		Apk:         a.Apk,
	}
}

// AppsHandler exposes the /api apps endpoints.
type AppsHandler struct {
	svc       *appsvc.Service
	maxUpload int64
}

// NewAppsHandler creates an AppsHandler.
func NewAppsHandler(svc *appsvc.Service, maxUpload int64) *AppsHandler {
	return &AppsHandler{svc: svc, maxUpload: maxUpload}
}

// List returns all apps as JSON.
func (h *AppsHandler) List(w http.ResponseWriter, r *http.Request) {
	apps, err := h.svc.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	out := make([]appDTO, 0, len(apps))
	for _, a := range apps {
		out = append(out, toDTO(a))
	}

	writeJSON(w, http.StatusOK, out)
}

// Upload accepts a multipart file, parses it and stores the app.
func (h *AppsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// net/http requires an explicit multipart parse before FormFile (gin did
	// it implicitly). MaxBytesReader caps the body at the configured limit.
	r.Body = http.MaxBytesReader(w, r.Body, h.maxUpload)
	if err := r.ParseMultipartForm(h.maxUpload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing file"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing file"})
		return
	}
	defer file.Close()

	if err := h.svc.UploadApk(r.Context(), file, header.Filename); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete removes the app identified by the :id path param.
func (h *AppsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

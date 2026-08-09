package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
func (h *AppsHandler) List(c *gin.Context) {
	apps, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]appDTO, 0, len(apps))
	for _, a := range apps {
		out = append(out, toDTO(a))
	}

	c.JSON(http.StatusOK, out)
}

// Upload accepts a multipart file, parses it and stores the app.
func (h *AppsHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxUpload)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open file"})
		return
	}
	defer f.Close()

	if err := h.svc.UploadApk(c.Request.Context(), f, file.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Delete removes the app identified by the :id path param.
func (h *AppsHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

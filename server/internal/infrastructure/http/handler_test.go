package http

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	appsvc "github.com/isdenmois/appdroid/server/internal/application/app"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/apkparser"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/apkstorage"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/http/views"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/repository"
)

// testApp returns the (handler bundle, data dir) used by the tests. It uses a
// real SQLite database and file storage isolated in a temp dir, so the tests
// exercise the whole stack except the network layer.
func testApp(t *testing.T) (*Handler, string) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	dataDir := t.TempDir()
	db, err := repository.Open(dataDir)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	svc := appsvc.NewService(
		repository.NewAppRepository(db),
		apkparser.NewParser(),
		apkstorage.NewStorage(dataDir),
	)

	renderer, err := views.NewRenderer()
	if err != nil {
		t.Fatalf("renderer: %v", err)
	}

	return &Handler{
		Apps:  NewAppsHandler(svc, 256*1024*1024),
		Files: NewFilesHandler(svc),
		Pages: NewPagesHandler(svc, renderer),
	}, dataDir
}

func newRequest(t *testing.T, router *gin.Engine, method, path string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "multipart/form-data")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestPing(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/api/ping", nil)

	// assert
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "pong" {
		t.Errorf("expected pong, got %q", w.Body.String())
	}
}

func TestListEmpty(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/api/list", nil)

	// assert
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "[]" {
		t.Errorf("expected empty list, got %q", w.Body.String())
	}
}

func TestUploadThenListThenDelete(t *testing.T) {
	// arrange
	h, dataDir := testApp(t)
	router := New(h)

	src := filepath.Join("..", "..", "..", "..", "data", "com.isdenmois.appdroid.apk")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("no apk fixture: %v", err)
	}

	// act: upload
	body, contentType := multipartBody(t, src)
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// assert: upload ok
	if w.Code != http.StatusOK {
		t.Fatalf("upload expected 200, got %d (%s)", w.Code, w.Body.String())
	}

	// act: list
	w = newRequest(t, router, http.MethodGet, "/api/list", nil)

	// assert: list has one app
	var apps []appDTO
	if err := json.Unmarshal(w.Body.Bytes(), &apps); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}
	if apps[0].AppID != "com.isdenmois.appdroid" || apps[0].Name != "AppDroid" {
		t.Errorf("unexpected app: %+v", apps[0])
	}

	// assert: file stored
	if _, err := os.Stat(filepath.Join(dataDir, apps[0].Apk)); err != nil {
		t.Errorf("expected stored apk file: %v", err)
	}

	// act: delete
	w = newRequest(t, router, http.MethodDelete, "/api/"+apps[0].ID, nil)

	// assert: delete ok and row gone
	if w.Code != http.StatusOK {
		t.Fatalf("delete expected 200, got %d", w.Code)
	}
	w = newRequest(t, router, http.MethodGet, "/api/list", nil)
	if w.Body.String() != "[]" {
		t.Errorf("expected empty list after delete, got %q", w.Body.String())
	}
}

func TestUploadInvalidFile(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)
	garbage := filepath.Join(t.TempDir(), "not-an-apk.apk")
	if err := os.WriteFile(garbage, []byte("garbage"), 0o644); err != nil {
		t.Fatalf("write garbage: %v", err)
	}

	// act
	body, contentType := multipartBody(t, garbage)
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// assert: uploading a non-APK fails cleanly
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestGetAppPageAndFile(t *testing.T) {
	// arrange
	h, dataDir := testApp(t)
	router := New(h)

	src := filepath.Join("..", "..", "..", "..", "data", "com.isdenmois.appdroid.apk")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("no apk fixture: %v", err)
	}
	body, contentType := multipartBody(t, src)
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("upload expected 200, got %d", w.Code)
	}

	// act
	w = newRequest(t, router, http.MethodGet, "/app/com.isdenmois.appdroid", nil)

	// assert: page renders the download link
	if w.Code != http.StatusOK {
		t.Fatalf("app page expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "com.isdenmois.appdroid-2.1.apk") {
		t.Errorf("expected download link in page, got %q", w.Body.String())
	}

	// act: file serving
	w = newRequest(t, router, http.MethodGet, "/apk/com.isdenmois.appdroid.apk/com.isdenmois.appdroid-2.1.apk", nil)

	// assert
	if w.Code != http.StatusOK {
		t.Fatalf("apk expected 200, got %d", w.Code)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "com.isdenmois.appdroid.apk")); err != nil {
		t.Errorf("expected stored file, got %v", err)
	}
}

func TestAppsPageHasObtainiumLinks(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	src := filepath.Join("..", "..", "..", "..", "data", "com.isdenmois.appdroid.apk")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("no apk fixture: %v", err)
	}
	body, contentType := multipartBody(t, src)
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("upload expected 200, got %d", w.Code)
	}

	// act
	req = httptest.NewRequest(http.MethodGet, "/apps", nil)
	req.Host = "192.168.1.5"
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// assert
	if w.Code != http.StatusOK {
		t.Fatalf("apps page expected 200, got %d", w.Code)
	}
	b := w.Body.String()
	if !strings.Contains(b, "obtainium://app/") {
		t.Errorf("expected obtainium links in page, got %q", b)
	}
	if !strings.Contains(b, "http://192.168.1.5/app/") {
		t.Errorf("expected http scheme for 192 host, got %q", b)
	}
}

func TestAppPageMissingReturns404(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/app/com.missing.app", nil)

	// assert
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestStaticIndexServed(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/", nil)

	// assert
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "App droid APK uploader") {
		t.Errorf("expected index.html content, got %q", w.Body.String())
	}
}

func TestStaticAssetsCacheControl(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/style.css", nil)

	// assert: subresources are cached for a day
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Cache-Control"); got != "public, max-age=86400" {
		t.Errorf("expected Cache-Control public, max-age=86400, got %q", got)
	}
}

func TestEntryPageCacheControl(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/", nil)

	// assert: the HTML document always revalidates
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Cache-Control"); got != "no-cache" {
		t.Errorf("expected Cache-Control no-cache, got %q", got)
	}
}

func TestPagesCacheControl(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/apps", nil)

	// assert: SSR pages always revalidate
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Cache-Control"); got != "no-cache" {
		t.Errorf("expected Cache-Control no-cache, got %q", got)
	}
}

func TestStaticJsModuleServed(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/home.js", nil)

	// assert: JS modules are cached for a day and typed as JS
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Cache-Control"); got != "public, max-age=86400" {
		t.Errorf("expected Cache-Control public, max-age=86400, got %q", got)
	}
	if got := w.Header().Get("Content-Type"); got != "text/javascript; charset=utf-8" {
		t.Errorf("expected Content-Type text/javascript; charset=utf-8, got %q", got)
	}
}

func TestStaticAssetsNoCacheInDevMode(t *testing.T) {
	// arrange: gin mode is a package global, so save the current mode, then
	// set debug mode after testApp (which pins TestMode).
	prev := gin.Mode()
	h, _ := testApp(t)
	gin.SetMode(gin.DebugMode)
	t.Cleanup(func() { gin.SetMode(prev) })

	router := New(h)

	// act: entry document and subresources
	for _, path := range []string{"/", "/style.css", "/home.js"} {
		w := newRequest(t, router, http.MethodGet, path, nil)

		// assert: nothing is cached in dev mode
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 for %s, got %d", path, w.Code)
		}
		if got := w.Header().Get("Cache-Control"); got != "no-cache" {
			t.Errorf("expected Cache-Control no-cache for %s, got %q", path, got)
		}
	}
}

func TestStaticMissingServes404(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act: an unknown asset and the dropped /static prefix both miss
	for _, path := range []string{"/nope.js", "/static/style.css"} {
		w := newRequest(t, router, http.MethodGet, path, nil)

		// assert
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 for %s, got %d", path, w.Code)
		}
	}
}

func TestApiHasNoCacheControl(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act
	w := newRequest(t, router, http.MethodGet, "/api/ping", nil)

	// assert: API responses are not tagged with a caching policy
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Cache-Control"); got != "" {
		t.Errorf("expected no Cache-Control header, got %q", got)
	}
}

// multipartBody builds a multipart request body uploading the file at src and
// returns the body and its Content-Type.
func multipartBody(t *testing.T, src string) (*bytes.Buffer, string) {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filepath.Base(src))
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if _, err := fw.Write(data); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}
	return &buf, w.FormDataContentType()
}

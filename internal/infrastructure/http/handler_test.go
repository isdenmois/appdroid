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

	"github.com/go-chi/chi/v5"

	appsvc "github.com/isdenmois/appdroid/internal/application/app"
	"github.com/isdenmois/appdroid/internal/infrastructure/apkparser"
	"github.com/isdenmois/appdroid/internal/infrastructure/apkstorage"
	"github.com/isdenmois/appdroid/internal/infrastructure/http/views"
	"github.com/isdenmois/appdroid/internal/infrastructure/repository"
)

// testApp returns the (handler bundle, data dir) used by the tests. It uses a
// real bbolt database and file storage isolated in a temp dir, so the tests
// exercise the whole stack except the network layer.
func testApp(t *testing.T) (*Handler, string) {
	t.Helper()

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
		// Any non-empty key enables auth in the tests; callers add the
		// X-API-Key header to the requests they send.
		APIKey: "test-api-key",
	}, dataDir
}

func newRequest(t *testing.T, router *chi.Mux, method, path string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "multipart/form-data")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// newAuthedRequest is like newRequest but carries a valid X-API-Key header, so
// it can reach routes guarded by the auth middleware.
func newAuthedRequest(t *testing.T, router *chi.Mux, method, path string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "multipart/form-data")
	}
	req.Header.Set("X-API-Key", "test-api-key")
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
	req.Header.Set("X-API-Key", "test-api-key")
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
	w = newAuthedRequest(t, router, http.MethodDelete, "/api/"+apps[0].ID, nil)

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
	req.Header.Set("X-API-Key", "test-api-key")
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
	req.Header.Set("X-API-Key", "test-api-key")
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
	req.Header.Set("X-API-Key", "test-api-key")
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

	// assets are only cached in release mode
	t.Setenv("SERVER_MODE", "release")
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

	// assets are only cached in release mode
	t.Setenv("SERVER_MODE", "release")
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
	// arrange: with SERVER_MODE unset (the test default) the server runs in
	// dev mode, so static assets must not be cached.
	h, _ := testApp(t)

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

// uploadAuthedRequest uploads the fixture at src with the given X-API-Key header
// value (or empty to omit it) and returns the recorder.
func uploadAuthedRequest(t *testing.T, router *chi.Mux, apiKey string, src string) *httptest.ResponseRecorder {
	t.Helper()

	body, contentType := multipartBody(t, src)
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", contentType)
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// multipartBytes builds a multipart request uploading the given bytes under
// field/fileName and returns the body and its Content-Type.
func multipartBytes(t *testing.T, field, fileName string, data []byte) (*bytes.Buffer, string) {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile(field, fileName)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write(data); err != nil {
		t.Fatalf("write data: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}
	return &buf, w.FormDataContentType()
}

func TestAuthUploadWithoutKeyReturns401(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act: upload any bytes with no X-API-Key header at all
	body, ct := multipartBytes(t, "file", "dummy.apk", []byte("not-an-apk"))
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", ct)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// assert
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestAuthUploadWithWrongKeyReturns401(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act: upload any bytes with the wrong key
	body, ct := multipartBytes(t, "file", "dummy.apk", []byte("not-an-apk"))
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", ct)
	req.Header.Set("X-API-Key", "not-the-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// assert
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestAuthUploadWithKeySucceeds(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	src := filepath.Join("..", "..", "..", "..", "data", "com.isdenmois.appdroid.apk")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("no apk fixture: %v", err)
	}

	// act: upload with the correct key
	w := uploadAuthedRequest(t, router, "test-api-key", src)

	// assert: same as an unguarded upload
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestAuthDeleteWithoutKeyReturns401(t *testing.T) {
	// arrange
	h, _ := testApp(t)
	router := New(h)

	// act: delete with no X-API-Key header
	w := newRequest(t, router, http.MethodDelete, "/api/does-not-matter", nil)

	// assert
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestAuthFailsClosedWithoutConfiguredKey(t *testing.T) {
	// arrange: a handler with no key configured at all.
	h, _ := testApp(t)
	h.APIKey = ""
	router := New(h)

	// act: an upload that would carry a valid key anyway.
	body, ct := multipartBytes(t, "file", "dummy.apk", []byte("not-an-apk"))
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", ct)
	req.Header.Set("X-API-Key", "test-api-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// assert: still rejected because no key is configured
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (%s)", w.Code, w.Body.String())
	}
}

// TestMutationRoutesAreRateLimited verifies that mutating API routes are
// rate-limited per client IP: once the configured limit is exceeded the
// middleware returns 429, while requests below the limit proceed.
func TestMutationRoutesAreRateLimited(t *testing.T) {
	// arrange
	h, _ := testApp(t)

	// Inject a small throttle limit; New applies middleware.Throttle(2) to the
	// mutating routes. All requests share one client IP so they share a bucket.
	h.throttleLimit = 2
	router := New(h)

	// act + assert: the first two mutations pass the limiter, the third is
	// rejected with 429 regardless of its authentication outcome.
	for i := 1; i <= 3; i++ {
		w := newAuthedRequest(t, router, "POST", "/api/upload", strings.NewReader(""))
		if i < 3 {
			if w.Code == http.StatusTooManyRequests {
				t.Fatalf("request %d: limiter blocked before the limit was reached", i)
			}
		} else if w.Code != http.StatusTooManyRequests {
			t.Fatalf("request 3: expected 429 after exceeding the limit, got %d", w.Code)
		}
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

package app

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	domainapp "github.com/isdenmois/appdroid/internal/domain/app"
)

// fakeRepo is an in-memory AppRepository for tests.
type fakeRepo struct {
	apps    map[string]domainapp.App
	adds    []domainapp.App
	deletes []string
}

func (f *fakeRepo) GetAll(_ context.Context) ([]domainapp.App, error) {
	out := make([]domainapp.App, 0, len(f.apps))
	for _, a := range f.apps {
		out = append(out, a)
	}
	return out, nil
}

func (f *fakeRepo) Get(_ context.Context, id string) (*domainapp.App, error) {
	if a, ok := f.apps[id]; ok {
		return &a, nil
	}
	return nil, nil
}

func (f *fakeRepo) Add(_ context.Context, a domainapp.App) error {
	f.adds = append(f.adds, a)
	f.apps[a.ID] = a
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, id string) error {
	f.deletes = append(f.deletes, id)
	delete(f.apps, id)
	return nil
}

// fakeParser returns the given metadata, optionally failing.
type fakeParser struct {
	md  domainapp.ApkMetadata
	err error
}

func (f *fakeParser) Parse(string) (domainapp.ApkMetadata, error) { return f.md, f.err }

// fakeStorage records the calls it receives and fails on demand.
type fakeStorage struct {
	savedTemp    string
	removedTemp  []string
	savedFiles   []string
	removedFiles []string
	saveErr      error
	openName     string
}

func (f *fakeStorage) SaveTemp(io.Reader) (string, error) { return "/tmp/fake", nil }
func (f *fakeStorage) RemoveTemp(path string) error {
	f.removedTemp = append(f.removedTemp, path)
	return nil
}
func (f *fakeStorage) SaveFile(src, name string) error {
	f.savedFiles = append(f.savedFiles, name)
	return f.saveErr
}
func (f *fakeStorage) Path(name string) string { return "/data/" + name }
func (f *fakeStorage) Remove(name string) error {
	f.removedFiles = append(f.removedFiles, name)
	return nil
}
func (f *fakeStorage) Open(name string) (io.ReadCloser, error) {
	f.openName = name
	return io.NopCloser(strings.NewReader("apk")), nil
}

func newFakeService() (*Service, *fakeRepo, *fakeStorage) {
	repo := &fakeRepo{apps: map[string]domainapp.App{}}
	storage := &fakeStorage{}
	return NewService(repo, &fakeParser{md: validMetadata()}, storage), repo, storage
}

func validMetadata() domainapp.ApkMetadata {
	return domainapp.ApkMetadata{
		AppID:       "com.example.app",
		Name:        "Example App",
		Version:     "42",
		VersionName: "1.2.3",
	}
}

func TestUploadApkStoresAppAndFile(t *testing.T) {
	// arrange
	svc, repo, storage := newFakeService()

	// act
	err := svc.UploadApk(context.Background(), strings.NewReader("apk"), "app.apk")

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.adds) != 1 {
		t.Fatalf("expected 1 stored app, got %d", len(repo.adds))
	}
	stored := repo.adds[0]
	if stored.ID != "com.example.app" {
		t.Errorf("expected id com.example.app, got %q", stored.ID)
	}
	if stored.Apk != "com.example.app.apk" {
		t.Errorf("expected apk com.example.app.apk, got %q", stored.Apk)
	}
	if len(storage.savedFiles) != 1 || storage.savedFiles[0] != "com.example.app.apk" {
		t.Errorf("expected saved file com.example.app.apk, got %v", storage.savedFiles)
	}
	if len(storage.removedTemp) != 1 {
		t.Errorf("expected temp cleanup, got %v", storage.removedTemp)
	}
}

func TestUploadApkRejectsMissingMetadata(t *testing.T) {
	// arrange
	svc, repo, _ := newFakeService()
	svc.parser = &fakeParser{md: domainapp.ApkMetadata{Name: "No Id"}}

	// act
	err := svc.UploadApk(context.Background(), strings.NewReader("apk"), "app.apk")

	// assert
	if !errors.Is(err, domainapp.ErrInvalidMetadata) {
		t.Fatalf("expected ErrInvalidMetadata, got %v", err)
	}
	if len(repo.adds) != 0 {
		t.Errorf("expected no stored app, got %d", len(repo.adds))
	}
}

func TestUploadApkParserErrorAborts(t *testing.T) {
	// arrange
	svc, repo, _ := newFakeService()
	svc.parser = &fakeParser{err: errors.New("boom")}

	// act
	err := svc.UploadApk(context.Background(), strings.NewReader("apk"), "app.apk")

	// assert
	if err == nil {
		t.Fatal("expected an error")
	}
	if len(repo.adds) != 0 {
		t.Errorf("expected no stored app, got %d", len(repo.adds))
	}
}

func TestGetReturnsNotFoundForMissingApp(t *testing.T) {
	// arrange
	svc, _, _ := newFakeService()

	// act
	_, err := svc.Get(context.Background(), "com.missing.app")

	// assert
	if !errors.Is(err, ErrAppNotFound) {
		t.Fatalf("expected ErrAppNotFound, got %v", err)
	}
}

func TestDeleteRemovesFileAndRow(t *testing.T) {
	// arrange
	svc, repo, storage := newFakeService()
	md := validMetadata()
	app, _ := domainapp.NewApp(md)
	repo.apps[app.ID] = *app

	// act
	err := svc.Delete(context.Background(), app.ID)

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(storage.removedFiles) != 1 || storage.removedFiles[0] != app.Apk {
		t.Errorf("expected removed file %q, got %v", app.Apk, storage.removedFiles)
	}
	if _, ok := repo.apps[app.ID]; ok {
		t.Errorf("expected app %q to be deleted", app.ID)
	}
}

func TestDeleteMissingAppIsNoop(t *testing.T) {
	// arrange
	svc, _, storage := newFakeService()

	// act
	err := svc.Delete(context.Background(), "com.missing.app")

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(storage.removedFiles) != 0 {
		t.Errorf("expected no file removed, got %v", storage.removedFiles)
	}
}

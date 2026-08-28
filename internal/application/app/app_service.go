// Package app contains the application (use-case) layer for the "apps"
// bounded context. It orchestrates the domain model and the infrastructure
// adapters without depending on any of them concretely.
package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	domainapp "github.com/isdenmois/appdroid/internal/domain/app"
)

// ErrAppNotFound is returned when the requested app does not exist.
var ErrAppNotFound = errors.New("app: not found")

// ApkParser extracts metadata from an APK file. It is an output port whose
// implementation lives in the infrastructure layer.
type ApkParser interface {
	// Parse reads the APK at path and returns its metadata.
	Parse(path string) (domainapp.ApkMetadata, error)
}

// ApkStorage persists APK files. It is an output port whose implementation
// lives in the infrastructure layer.
type ApkStorage interface {
	// SaveTemp stores the uploaded stream to a temporary file and returns its
	// path. The caller must remove the file via RemoveTemp.
	SaveTemp(r io.Reader) (string, error)
	// RemoveTemp deletes a file created by SaveTemp.
	RemoveTemp(path string) error
	// SaveFile copies the file at src into the data directory under name.
	SaveFile(src, name string) error
	// Path returns the data directory path of the stored apk name.
	Path(name string) string
	// Remove deletes the stored apk name from the data directory.
	Remove(name string) error
	// Open returns a read stream for the stored apk name.
	Open(name string) (io.ReadCloser, error)
}

// Service implements the use cases of the "apps" context.
type Service struct {
	repo    domainapp.AppRepository
	parser  ApkParser
	storage ApkStorage
}

// NewService creates a Service with its dependencies.
func NewService(repo domainapp.AppRepository, parser ApkParser, storage ApkStorage) *Service {
	return &Service{repo: repo, parser: parser, storage: storage}
}

// UploadApk validates the uploaded APK, extracts its metadata, stores the file
// under "<appId>.apk" and upserts the app record.
func (s *Service) UploadApk(ctx context.Context, file io.Reader, _ string) error {
	temp, err := s.storage.SaveTemp(file)
	if err != nil {
		return fmt.Errorf("save temp apk: %w", err)
	}
	defer s.storage.RemoveTemp(temp)

	md, err := s.parser.Parse(temp)
	if err != nil {
		return fmt.Errorf("parse apk: %w", err)
	}

	app, err := domainapp.NewApp(md)
	if err != nil {
		return fmt.Errorf("build app: %w", err)
	}

	if err := s.storage.SaveFile(temp, app.Apk); err != nil {
		return fmt.Errorf("save apk: %w", err)
	}

	if err := s.repo.Add(ctx, *app); err != nil {
		return fmt.Errorf("store app: %w", err)
	}

	return nil
}

// List returns all stored apps.
func (s *Service) List(ctx context.Context) ([]domainapp.App, error) {
	return s.repo.GetAll(ctx)
}

// Get returns a single app by id.
func (s *Service) Get(ctx context.Context, id string) (*domainapp.App, error) {
	app, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, ErrAppNotFound
	}
	return app, nil
}

// Delete removes the app record and its stored APK file. It is a no-op when
// the app does not exist.
func (s *Service) Delete(ctx context.Context, id string) error {
	app, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return nil
	}

	if err := s.storage.Remove(app.Apk); err != nil {
		return fmt.Errorf("remove apk: %w", err)
	}

	return s.repo.Delete(ctx, id)
}

// OpenFile returns a read stream for the stored APK of an existing app.
func (s *Service) OpenFile(ctx context.Context, name string) (io.ReadCloser, error) {
	return s.storage.Open(name)
}

// FilePath returns the data directory path of a stored apk name.
func (s *Service) FilePath(name string) string {
	return s.storage.Path(name)
}

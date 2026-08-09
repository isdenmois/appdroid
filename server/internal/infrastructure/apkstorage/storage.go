// Package apkstorage stores uploaded APK files in the data directory and
// provides temporary storage for uploads being processed.
package apkstorage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Storage persists APK files. It implements application.ApkStorage.
type Storage struct {
	// dataDir is where uploaded APKs live (default "./data").
	dataDir string
	// tmpDir is where uploads are staged before parsing (os temp dir).
	tmpDir string
}

// NewStorage creates a Storage writing permanent files into dataDir.
func NewStorage(dataDir string) *Storage {
	return &Storage{
		dataDir: dataDir,
		tmpDir:  filepath.Join(os.TempDir(), "apks"),
	}
}

// SaveTemp stores the stream to a temporary file and returns its path.
func (s *Storage) SaveTemp(r io.Reader) (string, error) {
	if err := os.MkdirAll(s.tmpDir, 0o755); err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	path := filepath.Join(s.tmpDir, newID())
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}

	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(path)
		return "", fmt.Errorf("write temp file: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("close temp file: %w", err)
	}

	return path, nil
}

// RemoveTemp deletes a file created by SaveTemp. It is a no-op when the file
// is not inside the temp dir or does not exist.
func (s *Storage) RemoveTemp(path string) error {
	if !filepath.HasPrefix(path, s.tmpDir) {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove temp apk: %w", err)
	}
	return nil
}

// SaveFile copies the file at src into the data directory under name.
func (s *Storage) SaveFile(src, name string) error {
	if err := os.MkdirAll(s.dataDir, 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer in.Close()

	out, err := os.Create(filepath.Join(s.dataDir, name))
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return fmt.Errorf("copy apk: %w", err)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("close destination: %w", err)
	}

	return nil
}

// Path returns the data directory path of the stored apk name.
func (s *Storage) Path(name string) string {
	return filepath.Join(s.dataDir, name)
}

// Remove deletes the stored apk name from the data directory.
func (s *Storage) Remove(name string) error {
	if err := os.Remove(filepath.Join(s.dataDir, name)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove apk: %w", err)
	}
	return nil
}

// Open returns a read stream for the stored apk name.
func (s *Storage) Open(name string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.dataDir, name))
}

// newID returns a random hex string used as the temporary file name
// (equivalent of the previous randomUUID()).
func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

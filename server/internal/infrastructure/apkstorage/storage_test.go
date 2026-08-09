package apkstorage

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveTempAndRemoveTemp(t *testing.T) {
	// arrange
	s := NewStorage(t.TempDir())
	content := "apk-bytes"

	// act
	path, err := s.SaveTemp(strings.NewReader(content))
	if err != nil {
		t.Fatalf("save temp: %v", err)
	}
	got, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read temp: %v", readErr)
	}
	removeErr := s.RemoveTemp(path)

	// assert
	if string(got) != content {
		t.Errorf("expected %q, got %q", content, string(got))
	}
	if removeErr != nil {
		t.Fatalf("remove temp: %v", removeErr)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected temp file to be removed, stat err=%v", err)
	}
}

func TestRemoveTempIgnoresForeignPaths(t *testing.T) {
	// arrange
	s := NewStorage(t.TempDir())
	foreign := filepath.Join(t.TempDir(), "other.tmp")
	if err := os.WriteFile(foreign, []byte("x"), 0o644); err != nil {
		t.Fatalf("write foreign: %v", err)
	}

	// act
	err := s.RemoveTemp(foreign)

	// assert
	if err != nil {
		t.Fatalf("remove temp: %v", err)
	}
	if _, statErr := os.Stat(foreign); statErr != nil {
		t.Errorf("expected foreign file to be untouched, stat err=%v", statErr)
	}
}

func TestSaveFileCopiesIntoDataDir(t *testing.T) {
	// arrange
	dataDir := t.TempDir()
	s := NewStorage(dataDir)
	src := filepath.Join(t.TempDir(), "src.apk")
	if err := os.WriteFile(src, []byte("content"), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}

	// act
	err := s.SaveFile(src, "com.example.app.apk")
	got, readErr := os.ReadFile(filepath.Join(dataDir, "com.example.app.apk"))

	// assert
	if err != nil {
		t.Fatalf("save file: %v", err)
	}
	if readErr != nil {
		t.Fatalf("read saved: %v", readErr)
	}
	if string(got) != "content" {
		t.Errorf("expected content, got %q", string(got))
	}
}

func TestPathAndOpen(t *testing.T) {
	// arrange
	dataDir := t.TempDir()
	s := NewStorage(dataDir)
	name := "com.example.app.apk"
	if err := os.WriteFile(filepath.Join(dataDir, name), []byte("apk"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// act
	path := s.Path(name)
	f, err := s.Open(name)
	defer f.Close()
	got, readErr := io.ReadAll(f)

	// assert
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if path != filepath.Join(dataDir, name) {
		t.Errorf("expected path %q, got %q", filepath.Join(dataDir, name), path)
	}
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}
	if string(got) != "apk" {
		t.Errorf("expected apk, got %q", string(got))
	}
}

func TestRemoveDeletesStoredFile(t *testing.T) {
	// arrange
	dataDir := t.TempDir()
	s := NewStorage(dataDir)
	name := "com.example.app.apk"
	if err := os.WriteFile(filepath.Join(dataDir, name), []byte("apk"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// act
	err := s.Remove(name)

	// assert
	if err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dataDir, name)); !os.IsNotExist(statErr) {
		t.Errorf("expected file to be removed, stat err=%v", statErr)
	}
}

func TestRemoveMissingFileIsNoop(t *testing.T) {
	// arrange
	s := NewStorage(t.TempDir())

	// act
	err := s.Remove("com.missing.app.apk")

	// assert
	if err != nil {
		t.Fatalf("remove: %v", err)
	}
}

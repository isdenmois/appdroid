package app

import (
	"errors"
	"testing"
)

func TestNewAppBuildsValidApp(t *testing.T) {
	// arrange
	md := ApkMetadata{
		AppID:       "com.example.app",
		Name:        "Example App",
		Version:     "42",
		VersionName: "1.2.3",
	}

	// act
	a, err := NewApp(md)

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID != "com.example.app" {
		t.Errorf("expected id com.example.app, got %q", a.ID)
	}
	if a.Type != TypeStatic {
		t.Errorf("expected type static, got %q", a.Type)
	}
	if a.Apk != "com.example.app.apk" {
		t.Errorf("expected apk com.example.app.apk, got %q", a.Apk)
	}
}

func TestNewAppRejectsMissingAppID(t *testing.T) {
	// arrange
	md := ApkMetadata{Name: "No Id"}

	// act
	_, err := NewApp(md)

	// assert
	if !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("expected ErrInvalidMetadata, got %v", err)
	}
}

func TestNewAppRejectsMissingName(t *testing.T) {
	// arrange
	md := ApkMetadata{AppID: "com.example.app"}

	// act
	_, err := NewApp(md)

	// assert
	if !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("expected ErrInvalidMetadata, got %v", err)
	}
}

package apkparser

import (
	"os"
	"path/filepath"
	"testing"
)

// fixture returns the path of an APK used as the parsing fixture. It looks for
// one of the APKs already present in the repo data directory.
func fixture(t *testing.T) string {
	t.Helper()

	for _, name := range []string{
		filepath.Join("..", "..", "..", "..", "data", "com.isdenmois.appdroid.apk"),
	} {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	t.Skip("no apk fixture available")
	return ""
}

func TestParseExtractsMetadata(t *testing.T) {
	// arrange
	p := NewParser()

	// act
	md, err := p.Parse(fixture(t))

	// assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if md.AppID != "com.isdenmois.appdroid" {
		t.Errorf("expected com.isdenmois.appdroid, got %q", md.AppID)
	}
	if md.Name != "AppDroid" {
		t.Errorf("expected AppDroid, got %q", md.Name)
	}
	if md.Version == "" {
		t.Error("expected a versionCode")
	}
	if md.VersionName == "" {
		t.Error("expected a versionName")
	}
}

func TestParseMissingFileReturnsError(t *testing.T) {
	// arrange
	p := NewParser()

	// act
	_, err := p.Parse(filepath.Join(t.TempDir(), "missing.apk"))

	// assert
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
}

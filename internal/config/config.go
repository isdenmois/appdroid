// Package config reads the runtime configuration from environment variables.
package config

import "os"

// Config holds the runtime settings of the server.
type Config struct {
	// Port is the HTTP listen port (default 3000).
	Port string
	// DataDir is the directory for the bbolt database and stored APKs
	// (default "./data").
	DataDir string
	// MaxUploadBytes caps the size of an uploaded APK (default 256 MiB).
	MaxUploadBytes int64
	// APIKey is the shared secret clients must send in the X-API-Key header to
	// perform mutating requests (POST/PATCH/DELETE). When empty the server
	// fails closed: every mutating request is rejected until a key is set.
	APIKey string
}

// Load reads the configuration from the environment.
func Load() Config {
	return Config{
		Port:           envOr("PORT", "3000"),
		DataDir:        envOr("DATA_DIR", "./data"),
		MaxUploadBytes: 256 * 1024 * 1024,
		APIKey:         os.Getenv("API_KEY"),
	}
}

// IsReleaseMode reports whether the server is running in release mode.
//
// Release mode is signalled by the SERVER_MODE environment variable (the old
// GIN_MODE variable is gone); any value other than "release" (including unset)
// means dev mode.
func IsReleaseMode() bool {
	return os.Getenv("SERVER_MODE") == "release"
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

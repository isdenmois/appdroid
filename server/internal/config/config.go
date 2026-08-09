// Package config reads the runtime configuration from environment variables.
package config

import "os"

// Config holds the runtime settings of the server.
type Config struct {
	// Port is the HTTP listen port (default 3000).
	Port string
	// DataDir is the directory for the SQLite database and stored APKs
	// (default "./data").
	DataDir string
	// MaxUploadBytes caps the size of an uploaded APK (default 256 MiB).
	MaxUploadBytes int64
}

// Load reads the configuration from the environment.
func Load() Config {
	return Config{
		Port:           envOr("PORT", "3000"),
		DataDir:        envOr("DATA_DIR", "./data"),
		MaxUploadBytes: 256 * 1024 * 1024,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

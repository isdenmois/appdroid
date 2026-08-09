// Package repository implements the AppRepository port on top of SQLite.
//
// The database file is opened with modernc.org/sqlite, a pure-Go driver that
// requires no CGO. The schema matches the one created by Drizzle in the
// previous Bun implementation, so an existing database is migrated in place.
package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS "apps" (
	"id" text PRIMARY KEY NOT NULL,
	"appId" text NOT NULL,
	"name" text NOT NULL,
	"version" text NOT NULL,
	"versionName" text NOT NULL,
	"type" text NOT NULL,
	"apk" text NOT NULL
);
`

// Open opens (creating if needed) and migrates the SQLite database located in
// dataDir. It returns a *sql.DB and a close function.
func Open(dataDir string) (*sql.DB, error) {
	path := filepath.Join(dataDir, "sqlite.db")

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate sqlite database: %w", err)
	}

	return db, nil
}

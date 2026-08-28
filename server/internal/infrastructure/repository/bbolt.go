// Package repository implements the AppRepository port on top of bbolt.
package repository

import (
	"fmt"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

// Open opens (creating if needed) the bbolt database located in dataDir.
//
// bbolt does not create the parent directory, so it is created up front. The
// "apps" bucket is created eagerly so the schema exists at open time.
func Open(dataDir string) (*bolt.DB, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	path := filepath.Join(dataDir, "appdroid.db")

	db, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open bbolt database: %w", err)
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("apps"))
		return err
	}); err != nil {
		db.Close()
		return nil, fmt.Errorf("create apps bucket: %w", err)
	}

	return db, nil
}

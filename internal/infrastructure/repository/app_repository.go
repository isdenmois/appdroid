package repository

import (
	"context"
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"

	domainapp "github.com/isdenmois/appdroid/internal/domain/app"
)

const appsBucket = "apps"

// AppRepository is the bbolt implementation of the domain AppRepository port.
//
// Methods accept a context.Context for interface compatibility, but it is not
// used: bbolt runs in-process and its transactions cannot be cancelled.
type AppRepository struct {
	db *bolt.DB
}

// NewAppRepository creates an AppRepository backed by db.
func NewAppRepository(db *bolt.DB) *AppRepository {
	return &AppRepository{db: db}
}

// encodeApp serializes an app to its bbolt value: a JSON document keyed by id.
func encodeApp(a domainapp.App) ([]byte, error) {
	return json.Marshal(a)
}

// decodeApp parses a bbolt value back into an app.
func decodeApp(b []byte) (domainapp.App, error) {
	var a domainapp.App
	err := json.Unmarshal(b, &a)
	return a, err
}

// GetAll returns all stored apps.
func (r *AppRepository) GetAll(ctx context.Context) ([]domainapp.App, error) {
	var apps []domainapp.App
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(appsBucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			a, err := decodeApp(v)
			if err != nil {
				return err
			}
			apps = append(apps, a)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list apps: %w", err)
	}
	return apps, nil
}

// Get returns a single app by id, or nil when it does not exist.
func (r *AppRepository) Get(ctx context.Context, id string) (*domainapp.App, error) {
	var app *domainapp.App
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(appsBucket))
		v := bucket.Get([]byte(id))
		if v == nil {
			return nil
		}
		a, err := decodeApp(v)
		if err != nil {
			return err
		}
		app = &a
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("get app %q: %w", id, err)
	}
	return app, nil
}

// Add inserts or updates an app by its id (bbolt Put overwrites).
func (r *AppRepository) Add(ctx context.Context, a domainapp.App) error {
	encoded, err := encodeApp(a)
	if err != nil {
		return fmt.Errorf("upsert app %q: %w", a.ID, err)
	}
	if err := r.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(appsBucket)).Put([]byte(a.ID), encoded)
	}); err != nil {
		return fmt.Errorf("upsert app %q: %w", a.ID, err)
	}
	return nil
}

// Delete removes an app by its id. A missing id is not an error.
func (r *AppRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(appsBucket)).Delete([]byte(id))
	}); err != nil {
		return fmt.Errorf("delete app %q: %w", id, err)
	}
	return nil
}

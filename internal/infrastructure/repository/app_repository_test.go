package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	domainapp "github.com/isdenmois/appdroid/internal/domain/app"
	bolt "go.etcd.io/bbolt"
)

// openTestDB opens an isolated bbolt database in a temp dir.
func openTestDB(t *testing.T) *AppRepository {
	t.Helper()

	db, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return NewAppRepository(db)
}

func sampleApp() domainapp.App {
	return domainapp.App{
		ID:          "com.example.app",
		AppID:       "com.example.app",
		Name:        "Example App",
		Version:     "42",
		VersionName: "1.2.3",
		Type:        domainapp.TypeStatic,
		Apk:         "com.example.app.apk",
	}
}

func TestRepositoryAddAndGet(t *testing.T) {
	// arrange
	repo := openTestDB(t)
	app := sampleApp()

	// act
	err := repo.Add(context.Background(), app)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	got, err := repo.Get(context.Background(), app.ID)

	// assert
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected an app, got nil")
	}
	if got.Name != app.Name || got.VersionName != app.VersionName {
		t.Errorf("expected %+v, got %+v", app, *got)
	}
}

func TestRepositoryGetMissingReturnsNil(t *testing.T) {
	// arrange
	repo := openTestDB(t)

	// act
	got, err := repo.Get(context.Background(), "com.missing.app")

	// assert
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %+v", *got)
	}
}

func TestRepositoryAddUpsertsByID(t *testing.T) {
	// arrange
	repo := openTestDB(t)
	app := sampleApp()
	if err := repo.Add(context.Background(), app); err != nil {
		t.Fatalf("add: %v", err)
	}

	// act
	updated := app
	updated.VersionName = "2.0.0"
	if err := repo.Add(context.Background(), updated); err != nil {
		t.Fatalf("add updated: %v", err)
	}
	got, err := repo.Get(context.Background(), app.ID)

	// assert
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil || got.VersionName != "2.0.0" {
		t.Fatalf("expected updated versionName, got %+v", got)
	}
	all, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("get all: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("expected 1 row after upsert, got %d", len(all))
	}
}

func TestRepositoryDelete(t *testing.T) {
	// arrange
	repo := openTestDB(t)
	app := sampleApp()
	if err := repo.Add(context.Background(), app); err != nil {
		t.Fatalf("add: %v", err)
	}

	// act
	err := repo.Delete(context.Background(), app.ID)
	got, getErr := repo.Get(context.Background(), app.ID)

	// assert
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if getErr != nil {
		t.Fatalf("get after delete: %v", getErr)
	}
	if got != nil {
		t.Errorf("expected nil after delete, got %+v", *got)
	}
}

func TestOpenCreatesDatabaseAndMigrates(t *testing.T) {
	// arrange
	dir := t.TempDir()

	// act
	db, err := Open(dir)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	// assert: the apps bucket exists and starts empty, and the database
	// file was created at the expected path.
	var count int
	var bucketExists bool
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(appsBucket))
		if b == nil {
			return nil
		}
		bucketExists = true
		for k, v := b.Cursor().First(); k != nil; k, v = b.Cursor().Next() {
			count++
			_ = v
		}
		return nil
	}); err != nil {
		t.Fatalf("view apps bucket: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "appdroid.db")); statErr != nil {
		t.Errorf("expected an appdroid.db file to be created: %v", statErr)
	}

	if !bucketExists {
		t.Error("expected the apps bucket to exist")
	}
	if count != 0 {
		t.Errorf("expected an empty apps bucket, got %d", count)
	}
}

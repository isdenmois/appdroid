package repository

import (
	"context"
	"database/sql"
	"fmt"

	domainapp "github.com/isdenmois/appdroid/server/internal/domain/app"
)

// AppRepository is the SQLite implementation of the domain AppRepository port.
type AppRepository struct {
	db *sql.DB
}

// NewAppRepository creates an AppRepository backed by db.
func NewAppRepository(db *sql.DB) *AppRepository {
	return &AppRepository{db: db}
}

// GetAll returns all stored apps.
func (r *AppRepository) GetAll(ctx context.Context) ([]domainapp.App, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, appId, name, version, versionName, type, apk FROM apps`)
	if err != nil {
		return nil, fmt.Errorf("list apps: %w", err)
	}
	defer rows.Close()

	var apps []domainapp.App
	for rows.Next() {
		var a domainapp.App
		if err := rows.Scan(&a.ID, &a.AppID, &a.Name, &a.Version, &a.VersionName, &a.Type, &a.Apk); err != nil {
			return nil, fmt.Errorf("scan app: %w", err)
		}
		apps = append(apps, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate apps: %w", err)
	}

	return apps, nil
}

// Get returns a single app by id, or nil when it does not exist.
func (r *AppRepository) Get(ctx context.Context, id string) (*domainapp.App, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, appId, name, version, versionName, type, apk FROM apps WHERE id = ?`, id)

	var a domainapp.App
	err := row.Scan(&a.ID, &a.AppID, &a.Name, &a.Version, &a.VersionName, &a.Type, &a.Apk)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get app %q: %w", id, err)
	}

	return &a, nil
}

// Add inserts or updates an app by its id.
func (r *AppRepository) Add(ctx context.Context, a domainapp.App) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO apps (id, appId, name, version, versionName, type, apk)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			appId = excluded.appId,
			name = excluded.name,
			version = excluded.version,
			versionName = excluded.versionName,
			type = excluded.type,
			apk = excluded.apk`,
		a.ID, a.AppID, a.Name, a.Version, a.VersionName, string(a.Type), a.Apk)
	if err != nil {
		return fmt.Errorf("upsert app %q: %w", a.ID, err)
	}
	return nil
}

// Delete removes an app by its id.
func (r *AppRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM apps WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete app %q: %w", id, err)
	}
	return nil
}

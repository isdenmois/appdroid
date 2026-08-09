package app

import "context"

// AppRepository is the port through which the application layer persists and
// reads App aggregates. Implementations live in the infrastructure layer.
type AppRepository interface {
	// GetAll returns all stored apps.
	GetAll(ctx context.Context) ([]App, error)
	// Get returns a single app by id, or nil when it does not exist.
	Get(ctx context.Context, id string) (*App, error)
	// Add inserts or updates (upsert) an app by its id.
	Add(ctx context.Context, a App) error
	// Delete removes an app by its id.
	Delete(ctx context.Context, id string) error
}

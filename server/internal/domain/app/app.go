// Package app contains the domain model for the "apps" bounded context.
//
// It is the innermost layer of the application and must not depend on any
// external packages: entities, value objects and business rules live here.
package app

import "errors"

// ErrInvalidMetadata is returned when an APK does not carry the data required
// to register an app (package id and application label).
var ErrInvalidMetadata = errors.New("app: invalid apk metadata")

// AppType is the kind of the stored app.
type AppType string

// TypeStatic is the only app type used by the service today.
const TypeStatic AppType = "static"

// App is the aggregate root representing an uploaded APK.
type App struct {
	ID          string
	AppID       string
	Name        string
	Version     string
	VersionName string
	Type        AppType
	Apk         string
}

// ApkMetadata is the set of fields extracted from an APK file by an
// infrastructure adapter (previously via the `aapt` binary).
type ApkMetadata struct {
	AppID       string
	Name        string
	Version     string
	VersionName string
}

// NewApp builds an App aggregate from extracted APK metadata.
//
// It enforces the domain invariants: the app id is used as the row id, the
// stored file name is "<appId>.apk" and the type is always "static".
func NewApp(md ApkMetadata) (*App, error) {
	if md.AppID == "" || md.Name == "" {
		return nil, ErrInvalidMetadata
	}

	return &App{
		ID:          md.AppID,
		AppID:       md.AppID,
		Name:        md.Name,
		Version:     md.Version,
		VersionName: md.VersionName,
		Type:        TypeStatic,
		Apk:         md.AppID + ".apk",
	}, nil
}

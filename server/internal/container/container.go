// Package container wires the application with sarulabs/di.
//
// It is the composition root: every definition is registered here and the
// dependencies are resolved through the container at build time.
package container

import (
	"go.etcd.io/bbolt"
	"net/http"

	"github.com/sarulabs/di/v2"

	appsvc "github.com/isdenmois/appdroid/server/internal/application/app"
	"github.com/isdenmois/appdroid/server/internal/config"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/apkparser"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/apkstorage"
	infrahttp "github.com/isdenmois/appdroid/server/internal/infrastructure/http"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/http/views"
	"github.com/isdenmois/appdroid/server/internal/infrastructure/repository"
)

const (
	// ConfigName is the container name of the runtime config.
	ConfigName = "config"
	// DBName is the container name of the bbolt database handle.
	DBName = "db"
	// AppServiceName is the container name of the application service.
	AppServiceName = "app-service"
	// HandlerName is the container name of the delivery layer bundle.
	HandlerName = "handler"

	repoName    = "app-repository"
	parserName  = "apk-parser"
	storageName = "apk-storage"
	viewsName   = "views"
)

// New builds the application container with every registered definition.
func New(cfg config.Config) (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return di.Container{}, err
	}

	configDef := di.Def{
		Name:  ConfigName,
		Build: func(di.Container) (interface{}, error) { return cfg, nil },
	}

	dbDef := di.Def{
		Name:  DBName,
		Build: func(di.Container) (interface{}, error) { return repository.Open(cfg.DataDir) },
		Close: func(obj interface{}) error {
			if db, ok := obj.(*bbolt.DB); ok {
				return db.Close()
			}
			return nil
		},
	}

	repoDef := di.Def{
		Name: repoName,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewAppRepository(ctn.Get(DBName).(*bbolt.DB)), nil
		},
	}

	parserDef := di.Def{
		Name:  parserName,
		Build: func(di.Container) (interface{}, error) { return apkparser.NewParser(), nil },
	}

	storageDef := di.Def{
		Name: storageName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(ConfigName).(config.Config)
			return apkstorage.NewStorage(cfg.DataDir), nil
		},
	}

	serviceDef := di.Def{
		Name: AppServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return appsvc.NewService(
				ctn.Get(repoName).(*repository.AppRepository),
				ctn.Get(parserName).(*apkparser.Parser),
				ctn.Get(storageName).(*apkstorage.Storage),
			), nil
		},
	}

	viewsDef := di.Def{
		Name:  viewsName,
		Build: func(di.Container) (interface{}, error) { return views.NewRenderer() },
	}

	handlerDef := di.Def{
		Name: HandlerName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(ConfigName).(config.Config)
			svc := ctn.Get(AppServiceName).(*appsvc.Service)

			apps := infrahttp.NewAppsHandler(svc, cfg.MaxUploadBytes)
			files := infrahttp.NewFilesHandler(svc)
			pages := infrahttp.NewPagesHandler(svc, ctn.Get(viewsName).(*views.Renderer))

			return &infrahttp.Handler{Apps: apps, Files: files, Pages: pages, APIKey: cfg.APIKey}, nil
		},
	}

	if err := builder.Add(configDef, dbDef, repoDef, parserDef, storageDef, serviceDef, viewsDef, handlerDef); err != nil {
		return di.Container{}, err
	}

	return builder.Build(), nil
}

// Router returns the configured chi router from the container.
func Router(ctn di.Container) http.Handler {
	return infrahttp.New(ctn.Get(HandlerName).(*infrahttp.Handler))
}

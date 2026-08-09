// Command server is the composition root of the AppDroid server.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isdenmois/appdroid/server/internal/config"
	"github.com/isdenmois/appdroid/server/internal/container"
)

func main() {
	cfg := config.Load()

	ctn, err := container.New(cfg)
	if err != nil {
		log.Fatalf("build container: %v", err)
	}
	defer ctn.Delete()

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: container.Router(ctn),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

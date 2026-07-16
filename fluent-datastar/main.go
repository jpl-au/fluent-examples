// Fluent-Datastar - a small example application demonstrating the
// fluent-datastar bindings. One server serves the same demo page twice:
// / drives backend endpoints built with the fluent-datastar SSE
// generator, and /sdk drives identical endpoints built with the
// official Datastar Go SDK. Both speak the same wire protocol, so the
// pages differ only in which implementation answers their buttons. Run
// it and visit http://localhost:8080 to see it in action.
package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jpl-au/fluent-examples/fluent-datastar/handler"
)

//go:embed static
var staticEmbed embed.FS

func main() {
	staticFS, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		slog.Error("failed to open embedded static assets", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	mux.HandleFunc("GET /{$}", handler.Home)
	mux.HandleFunc("POST /increment", handler.Increment)
	mux.HandleFunc("GET /clock", handler.Clock)
	mux.HandleFunc("GET /sdk", handler.SDKHome)
	mux.HandleFunc("POST /sdk/increment", handler.SDKIncrement)
	mux.HandleFunc("GET /sdk/clock", handler.SDKClock)

	srv := &http.Server{Addr: ":8080", Handler: mux}

	// Start the server in a goroutine so we can listen for shutdown
	// signals without blocking.
	go func() {
		slog.Info("fluent-datastar example listening on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// First signal starts graceful shutdown; the default handler takes
	// over for a second signal so it kills immediately.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	slog.Info("shutting down")
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdown); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}
	slog.Info("stopped")
}

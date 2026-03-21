// Contact Manager - a Fluent-JIT example application demonstrating
// JIT-optimised server-rendered HTML with plain HTTP, built with
// reusable components and the chain router. Run it and visit
// http://localhost:8080 to see it in action.
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

	"github.com/jpl-au/fluent-examples/fluent-jit/server"
)

//go:embed static
var staticEmbed embed.FS

func main() {
	staticFS, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		slog.Error("failed to open embedded static assets", "error", err)
		os.Exit(1)
	}

	handler := server.New(staticFS)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Start the server in a goroutine so we can listen for shutdown
	// signals without blocking.
	go func() {
		slog.Info("fluent-jit contact manager listening on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// First signal starts graceful shutdown. Second signal kills
	// immediately. NotifyContext stops listening after the first
	// signal so the default handler (process exit) takes over.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop() // Unregister so a second signal kills immediately.

	slog.Info("shutting down")
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdown); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}
	slog.Info("stopped")
}

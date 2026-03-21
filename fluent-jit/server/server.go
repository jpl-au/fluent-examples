// Package server creates and configures the HTTP server for the
// contact manager example.
package server

import (
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/jpl-au/chain"

	"github.com/jpl-au/fluent-examples/fluent-jit/routes"
)

// New creates a chain.Mux with middleware, static file serving, and
// all application routes registered. The staticFS should be the
// embedded static assets directory (already sub-ed to remove the
// "static" prefix).
func New(staticFS fs.FS) http.Handler {
	mux := chain.New()
	mux.Use(requestLogging)

	// Serve embedded static assets at /static/.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	routes.Register(mux)
	return mux
}

// requestLogging is simple middleware that logs each request with
// method, path, status, and duration using slog.
func requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		status := 0
		if rw, ok := w.(chain.ResponseWriter); ok {
			status = rw.Status()
		}

		slog.Info("http",
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"duration", time.Since(start),
		)
	})
}

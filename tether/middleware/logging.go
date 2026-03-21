// Package middleware provides example middleware for the feature
// explorer.
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jpl-au/chain"
	tether "github.com/jpl-au/tether"
)

// RequestLogging is standard HTTP middleware that logs every request
// with method, path, status code, and duration. It relies on
// chain.ResponseWriter to capture the status after the handler runs.
func RequestLogging(next http.Handler) http.Handler {
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

// Logging wraps the event handler to log each incoming event
// before passing it down the chain. Useful during development to
// trace event flow without adding slog calls in every handler.
func Logging[S any](next tether.HandleFunc[S]) tether.HandleFunc[S] {
	return func(sess tether.Session, s S, ev tether.Event) S {
		args := []any{
			"action", ev.Action,
			"type", ev.Type,
		}
		for k, v := range ev.Data {
			args = append(args, "data."+k, v)
		}
		if id := sess.ID(); id != "" {
			args = append(args, "session", id)
		}
		slog.Info("event", args...)
		return next(sess, s, ev)
	}
}

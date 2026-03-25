// Package store provides filesystem-backed implementations of
// [tether.SessionStore] and [tether.DiffStore] for the example
// application. Both write files to a directory on disk - simple,
// dependency-free, and easy to inspect during development.
//
// This package exists to demonstrate how to implement the store
// interfaces yourself. For a ready-made filesystem store, use
// github.com/jpl-au/tether-store/fs instead.
//
// Production applications would typically use Redis, SQLite, or
// another external store suited to their deployment.
package store

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// shortID returns the first six characters of an ID for logging,
// or the full ID if it is shorter than six characters.
func shortID(id string) string {
	if len(id) < 6 {
		return id
	}
	return id[:6]
}

// FileSessionStore persists session state to the filesystem. Each
// session is stored as a single file named by its ID. The TTL
// parameter from Save is logged but not enforced - the framework
// calls Delete on reconnect and destroy, so orphaned files are the
// only case where TTL would matter. A production store (Redis,
// SQL) would use TTL for automatic expiry.
type FileSessionStore struct {
	dir string
}

// NewFileSessionStore creates a FileSessionStore that writes session
// files to dir. The directory is created if it does not exist.
func NewFileSessionStore(dir string) *FileSessionStore {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		slog.Warn("session store: mkdir failed", "dir", dir, "error", err)
	}
	return &FileSessionStore{dir: dir}
}

// path builds the on-disk filename for a session by ID.
func (s *FileSessionStore) path(id string) string {
	return filepath.Join(s.dir, id+".session")
}

// Save writes session data to a file. The TTL is logged for
// observability but not enforced by the filesystem store.
func (s *FileSessionStore) Save(_ context.Context, id string, data []byte, ttl time.Duration) error {
	slog.Debug("session store: save", "id", shortID(id), "bytes", len(data), "ttl", ttl)
	return os.WriteFile(s.path(id), data, 0o600)
}

// Load reads session data from a file. Returns (nil, nil) if the
// file does not exist - the framework treats this as "no session
// to restore".
func (s *FileSessionStore) Load(_ context.Context, id string) ([]byte, error) {
	data, err := os.ReadFile(s.path(id))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err == nil {
		slog.Debug("session store: load", "id", shortID(id), "bytes", len(data))
	}
	return data, err
}

// Delete removes a session file. Returns nil if the file does not
// exist.
func (s *FileSessionStore) Delete(_ context.Context, id string) error {
	slog.Debug("session store: delete", "id", shortID(id))
	err := os.Remove(s.path(id))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// FileDiffStore persists differ snapshots to the filesystem. Each
// session's snapshots are stored as a single file named by its ID.
// This offloads snapshot data from Go memory during the reconnect
// window - a memory optimisation, not a recovery mechanism.
type FileDiffStore struct {
	dir string
}

// NewFileDiffStore creates a FileDiffStore that writes snapshot
// files to dir. The directory is created if it does not exist.
func NewFileDiffStore(dir string) *FileDiffStore {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		slog.Warn("diff store: mkdir failed", "dir", dir, "error", err)
	}
	return &FileDiffStore{dir: dir}
}

// path builds the on-disk filename for a diff snapshot by ID.
func (s *FileDiffStore) path(id string) string {
	return filepath.Join(s.dir, id+".diff")
}

// Save writes differ snapshot data to a file.
func (s *FileDiffStore) Save(_ context.Context, id string, data []byte) error {
	slog.Debug("diff store: save", "id", shortID(id), "bytes", len(data))
	return os.WriteFile(s.path(id), data, 0o600)
}

// Load reads differ snapshot data from a file. Returns (nil, nil)
// if the file does not exist.
func (s *FileDiffStore) Load(_ context.Context, id string) ([]byte, error) {
	data, err := os.ReadFile(s.path(id))
	if os.IsNotExist(err) {
		return nil, nil
	}
	return data, err
}

// Delete removes a differ snapshot file. Returns nil if the file
// does not exist.
func (s *FileDiffStore) Delete(_ context.Context, id string) error {
	slog.Debug("diff store: delete", "id", shortID(id))
	err := os.Remove(s.path(id))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

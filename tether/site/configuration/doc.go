// Package configuration demonstrates production-ready handler settings:
// tether.Timeouts (idle, reconnect, heartbeat), tether.Limits
// (max sessions, buffer size, event bytes), tether.Security
// (allowed origins, session binding), ws.Compression, and session
// persistence (SessionStore, DiffStore, OnRestore). Each configured
// value is rendered in an informational card so developers can see
// the available knobs and sensible defaults.
package configuration

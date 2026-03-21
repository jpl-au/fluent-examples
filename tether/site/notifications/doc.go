// Package notifications demonstrates server-initiated notification
// side effects: Toast, Flash, Announce, and Signal. All demos are
// pure side effects - the server pushes to the client without
// changing session state. Requires a persistent connection (WebSocket)
// because the server initiates communication at any time.
package notifications

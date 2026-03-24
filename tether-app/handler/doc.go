// Package handler provides the single tether.Handler for the kanban
// board application. It uses mode.Both for WebSocket with automatic
// SSE+POST failover. All board state lives in a shared store; per-session
// state tracks only the user's name, current view, and selected card.
package handler

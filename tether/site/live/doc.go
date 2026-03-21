// Package live demonstrates real-time features: uptime ticker via
// sess.Go, activity feeds via tether.Bus, online count via
// tether.Value, Group operations (Broadcast, BroadcastOthers, All,
// OnJoin/OnLeave), SetTitle, State() in background goroutines, and
// Close(). Includes WebSocket (full set) and SSE (subset) handlers.
package live

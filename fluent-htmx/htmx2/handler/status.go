package handler

import "fmt"

// setStatus returns the inline handler that writes a connection
// state into the status span with the given id. The SSE and
// WebSocket pages bind it to the htmx lifecycle events with HxOn.
func setStatus(id, text, class string) string {
	return fmt.Sprintf("var s=document.getElementById(%q);s.textContent=%q;s.className=%q", id, text, class)
}

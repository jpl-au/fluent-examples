// The /sdk endpoints. These mirror Increment and Clock in handler.go
// but speak the Datastar SSE protocol through the official Go SDK
// (datastar.NewSSE, PatchElements, MarshalAndPatchSignals) instead of
// the fluent-datastar generator. The element markup is still produced
// by fluent and handed to the SDK as a string, and the wire events are
// identical, so the same client page works against either half.
package handler

import (
	"log/slog"
	"net/http"
	"time"

	datastar "github.com/starfederation/datastar-go/datastar"
)

// SDKHome renders the same demo page wired to the SDK server half.
func SDKHome(w http.ResponseWriter, _ *http.Request) {
	render(w, "Official SDK",
		"This page drives endpoints built with the official Datastar Go SDK.",
		"/sdk/increment", "/sdk/clock",
		"/", "the same app driven by the fluent-datastar generator")
}

// SDKIncrement is the SDK-half round-trip endpoint, the counterpart of
// Increment in handler.go.
func SDKIncrement(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Count int `json:"count"`
	}
	if err := datastar.ReadSignals(r, &in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("increment: signals received", "source", "sdk", "count", in.Count)

	sse := datastar.NewSSE(w, r)

	if err := sse.PatchElements(string(resultPanel(in.Count).RenderBytes())); err != nil {
		slog.Error("patch elements", "error", err)
		return
	}
	if err := sse.MarshalAndPatchSignals(map[string]int{"serverCount": in.Count}); err != nil {
		slog.Error("patch signals", "error", err)
	}
}

// SDKClock is the SDK-half streaming endpoint, the counterpart of Clock
// in handler.go.
func SDKClock(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	slog.Info("clock: countdown started", "source", "sdk", "remote", r.RemoteAddr)

	ctx := r.Context()
	for n := countdownFrom; n >= 0; n-- {
		if err := sse.PatchElements(string(clockPanel(n).RenderBytes())); err != nil {
			return
		}
		if err := sse.MarshalAndPatchSignals(map[string]int{"tick": n}); err != nil {
			return
		}

		if n == 0 {
			slog.Info("clock: liftoff", "source", "sdk")
			return
		}
		select {
		case <-ctx.Done():
			slog.Info("clock: client disconnected", "source", "sdk", "remote", r.RemoteAddr)
			return
		case <-time.After(time.Second):
		}
	}
}

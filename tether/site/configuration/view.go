package configuration

import (
	"strconv"
	"strings"

	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/components/composite/configtable"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the configuration page with demo cards showing
// the Timeouts, Limits, Security, Compression, and Persistence
// configuration fields and their current values.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Timeouts",
			"Duration-based settings that control session lifecycle, reconnection, and transport keep-alive. Zero values use framework defaults.",
			"tether.Timeouts", panel.WS|panel.SSE,
			configtable.New(
				configtable.Row("Idle", configuredTimeouts.Idle.String()),
				configtable.Row("MaxLifetime", configuredTimeouts.MaxLifetime.String()),
				configtable.Row("Reconnect", configuredTimeouts.Reconnect.String()),
				configtable.Row("Pending", configuredTimeouts.Pending.String()),
				configtable.Row("ShutdownGrace", configuredTimeouts.ShutdownGrace.String()),
				configtable.Row("Heartbeat", configuredTimeouts.Heartbeat.String()),
				configtable.Row("Retry", configuredTimeouts.Retry.String()),
				configtable.Row("MaxRetry", configuredTimeouts.MaxRetry.String()),
				configtable.Row("BackoffMultiplier", strconv.FormatFloat(configuredTimeouts.BackoffMultiplier, 'f', 1, 64)),
				configtable.Row("DisableJitter", strconv.FormatBool(configuredTimeouts.DisableJitter)),
			),
		),

		panel.Card(
			"Limits",
			"Capacity constraints that protect against resource exhaustion. MaxSessions and MaxPending guard against flooding; CmdBufferSize tunes the per-session command channel.",
			"tether.Limits", panel.AllTransports,
			configtable.New(
				configtable.Row("MaxSessions", strconv.Itoa(configuredLimits.MaxSessions)),
				configtable.Row("MaxPending", strconv.Itoa(configuredLimits.MaxPending)),
				configtable.Row("CmdBufferSize", strconv.Itoa(configuredLimits.CmdBufferSize)),
				configtable.Row("MaxEventBytes", formatBytes(configuredLimits.MaxEventBytes)),
			),
		),

		panel.Card(
			"Security",
			"Cross-origin protection uses Go 1.25's http.CrossOriginProtection. Safe methods (GET, HEAD) are always allowed. State-changing requests are checked via Sec-Fetch-Site and Origin headers. Session binding verifies the User-Agent on reconnect to detect stolen session IDs.",
			"tether.Security", panel.AllTransports,
			configtable.New(
				configtable.Row("TrustedOrigins", strings.Join(configuredSecurity.TrustedOrigins, ", ")),
				configtable.Row("DisableSessionBinding", strconv.FormatBool(configuredSecurity.DisableSessionBinding)),
			),
		),
		panel.Card(
			"WebSocket Compression",
			"Per-message deflate (RFC 7692) is enabled by default. Browsers negotiate the extension transparently during the handshake. ContextTakeover retains the compression dictionary across messages for better ratios on repetitive HTML, at the cost of ~4 KB per connection.",
			"ws.Compression", panel.WS,
			configtable.New(
				configtable.Row("Disabled", strconv.FormatBool(configuredCompression.Disabled)),
				configtable.Row("Level", compressionLevelName(configuredCompression.Level)),
				configtable.Row("Threshold", strconv.Itoa(configuredCompression.Threshold)+" B"),
				configtable.Row("ContextTakeover", strconv.FormatBool(configuredCompression.ContextTakeover)),
			),
		),
		panel.Card(
			"Session Persistence",
			"SessionStore persists state to disk on disconnect and graceful shutdown, enabling crash recovery. DiffStore offloads differ snapshots during the reconnect window to free memory. OnRestore fires instead of OnConnect for recovered sessions - use it to rejoin groups, restart timers, or re-subscribe to buses.",
			"SessionStore · DiffStore · OnRestore", panel.WS|panel.SSE,
			configtable.New(
				configtable.Row("SessionStore", "FileSessionStore (/tmp/tether-sessions)"),
				configtable.Row("DiffStore", "FileDiffStore (/tmp/tether-diffs)"),
				configtable.Row("OnRestore", "Rejoins presence tracking, logs recovery"),
			),
		),
		panel.Card(
			"Page View Counter",
			"Bus.Emit inside OnNavigate fires during the initial GET (pre-warming) because CaptureSession.enqueue runs synchronously. The global subscriber counts every page view - including the very first render before any WebSocket connects.",
			"Bus.Emit · OnNavigate", panel.AllTransports,
			configtable.New(
				configtable.Row("Total Page Views", strconv.FormatInt(s.PageViews, 10)),
			),
		),
	)
}

// compressionLevelName returns a human-readable name for the level.
func compressionLevelName(level ws.CompressionLevel) string {
	switch level {
	case ws.CompressionFastest:
		return "Fastest (1)"
	case ws.CompressionBalanced:
		return "Balanced (6)"
	case ws.CompressionSmallest:
		return "Smallest (9)"
	default:
		return strconv.Itoa(int(level))
	}
}

// formatBytes returns a human-readable representation of a byte count
// (e.g. "128 KB"). Falls back to the raw byte count for values that
// are not cleanly divisible.
func formatBytes(b int64) string {
	switch {
	case b >= 1<<20 && b%(1<<20) == 0:
		return strconv.FormatInt(b/(1<<20), 10) + " MB"
	case b >= 1<<10 && b%(1<<10) == 0:
		return strconv.FormatInt(b/(1<<10), 10) + " KB"
	default:
		return strconv.FormatInt(b, 10) + " B"
	}
}

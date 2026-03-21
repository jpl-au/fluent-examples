package page

import (
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// LifecycleRender builds the PWA lifecycle events page, demonstrating
// online/offline connectivity detection and app installation events.
func LifecycleRender(s state.State) node.Node {
	return cpage.New(
		panel.Card(
			"Connectivity Events",
			"When the browser detects a network change, the client sends an online or offline event automatically - no bind helper needed. Try toggling your network in DevTools to see the status update. Note: when truly offline the WebSocket drops, so the offline event queues until the connection restores; the server learns about the offline period from the reconnection sequence.",
			"event.Online · event.Offline", panel.WS|panel.SSE,
			layout.Container(
				bind.Apply(panel.SignalSuccess("Online - connected to server"), bind.BindShow("sw.online")),
				bind.Apply(panel.SignalMuted("Offline - waiting for reconnection"), bind.BindHide("sw.online")),
			),
		),

		panel.Card(
			"App Installed",
			"When the user installs the PWA via the browser's install prompt, an appinstalled event fires automatically. The server can react - unlock features, update a user record, or show a thank-you message. This demo shows a confirmation when the event arrives.",
			"event.AppInstalled", panel.WS|panel.SSE,
			appInstalledResult(s.Lifecycle.Installed),
		),
	)
}

// LifecycleHandle processes PWA lifecycle events. These events arrive
// automatically from the client - no bind helper is needed.
func LifecycleHandle(sess tether.Session, s state.State, ev tether.Event) state.State {
	switch ev.Action {
	case "online":
		sess.Signal("sw.online", true)
	case "offline":
		sess.Signal("sw.online", false)
	case "appinstalled":
		s.Lifecycle.Installed = true
	}
	return s
}

// appInstalledResult renders a success message or an install prompt
// depending on whether the PWA installation event has fired.
func appInstalledResult(installed bool) node.Node {
	if installed {
		return p.Text("PWA installed - thank you!").Dynamic("app-installed")
	}
	return hint.Text("Install the app from your browser's address bar to see this event fire.").Dynamic("app-installed")
}

package page

import (
	"errors"

	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
	"github.com/jpl-au/tether/push"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// PushRender builds the push notifications page, demonstrating
// VAPID-authenticated Web Push subscription and notification delivery.
func PushRender(_ state.State) node.Node {
	return cpage.New(
		panel.Card(
			"Subscribe to Push",
			"Click the button to request notification permission and subscribe via the service worker's PushManager. The browser shows a permission dialog - this requires a genuine user gesture, which is why bind.PushSubscribe is attached to a button click.",
			"bind.PushSubscribe", panel.WS|panel.SSE,
			pushSubscribeSection(),
		),

		panel.Card(
			"Send Test Push",
			"Once subscribed, click to send a push notification from the server. The notification appears even if this tab is in the background or closed. The server uses VAPID keys generated on startup for authentication.",
			"sess.Push · push.Notification", panel.WS|panel.SSE,
			button.PrimaryAction("Send Test Push", "push.send",
				bind.BindShow("push.available"),
			),
		),

		panel.Card(
			"Rich Notification",
			"Push notifications support a title, body, icon, badge, URL, tag for grouping, and up to two action buttons. Click to send a notification with action buttons that navigate to specific pages.",
			"push.NotificationAction", panel.WS|panel.SSE,
			button.PrimaryAction("Send Rich Push", "push.rich",
				bind.BindShow("push.available"),
			),
		),
	)
}

// pushSubscribeSection renders the push subscription button when
// VAPID keys are available, or a fallback message when they are not.
// Visibility is controlled by signal bindings so the client toggles
// the right branch without a server round-trip.
func pushSubscribeSection() node.Node {
	return layout.Container(
		bind.Apply(
			layout.Row(
				button.Primary("Enable Push Notifications", bind.PushSubscribe()),
			),
			bind.BindShow("push.available"),
		),
		bind.Apply(
			hint.Text("Push notifications are unavailable (VAPID key generation failed)."),
			bind.BindShow("push.unavailable"),
		),
	)
}

// PushHandle processes events on the push page, sending a test
// notification when the user clicks the button.
func PushHandle(sess tether.Session, s state.State, ev tether.Event) state.State {
	switch ev.Action {
	case "push.send":
		if err := sess.Push(push.Notification{
			Title: "Test Push",
			Body:  "Hello from the server!",
		}); err != nil {
			pushError(sess, err)
		}
	case "push.rich":
		if err := sess.Push(push.Notification{
			Title: "New Activity",
			Body:  "Someone joined the session. Click to view.",
			Tag:   "activity",
			Actions: []push.NotificationAction{
				{Action: "view", Title: "View", URL: "/sw/"},
				{Action: "dismiss", Title: "Dismiss"},
			},
		}); err != nil {
			pushError(sess, err)
		}
	}
	return s
}

// pushError handles push delivery failures. Expired subscriptions
// (HTTP 410 from the push service) get a specific message prompting
// re-subscription; other errors surface the raw message.
func pushError(sess tether.Session, err error) {
	if errors.Is(err, push.ErrSubscriptionExpired) {
		sess.Toast("Subscription expired - please re-subscribe.")
		return
	}
	sess.Toast("Push failed: " + err.Error())
}

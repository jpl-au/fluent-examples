// Package feed provides styled list feeds for activity logs and
// message streams.
package feed

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
)

// Activity creates a styled activity feed from the given items.
func Activity(items ...node.Node) node.Node {
	return div.New(items...).Class("activity-feed")
}

// ActivityItem creates an activity feed entry with user, action, and timestamp.
func ActivityItem(user, action, time string) node.Node {
	return div.New(
		div.New(
			span.Text(user).Class("activity-user"),
			span.Text(time).Class("activity-time"),
		).Class("activity-header"),
		span.Text(action).Class("activity-action"),
	).Class("activity-item")
}

// ActivityText creates a simple text-only activity item.
func ActivityText(text string) node.Node {
	return div.New().Class("activity-item").Text(text)
}

// Messages creates a styled message feed from the given items.
func Messages(items ...node.Node) node.Node {
	return ul.New(items...).Class("msg-list")
}

// MessageItem creates a message feed entry with user and text.
func MessageItem(user, text string) node.Node {
	return li.New(
		span.Text(user).Class("msg-user"),
		span.Text(text).Class("msg-text"),
	).Class("msg-item")
}

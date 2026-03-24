// Package diagnostic provides styled elements for displaying
// diagnostic events and the diagnostic kind reference list.
package diagnostic

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
)

// List creates a container for diagnostic reference cards.
func List(items ...node.Node) node.Node {
	return div.New(items...).Class("diag-list")
}

// Item renders a single diagnostic reference card with kind label,
// description, and an optional trigger element (button, link, or text).
func Item(kind, desc string, trigger node.Node) node.Node {
	children := []node.Node{
		span.Text(kind).Class("diag-kind"),
		p.Text(desc).Class("diag-desc"),
	}
	if trigger != nil {
		children = append(children, trigger)
	}
	return div.New(children...).Class("diag-item")
}

// Trigger renders plain-text instructions for triggering a diagnostic.
func Trigger(s string) node.Node {
	return p.Text(s).Class("diag-trigger")
}

// EventList creates a container for live diagnostic events. Returns
// the concrete element so callers can chain .Dynamic().
func EventList(items ...node.Node) *ul.Element {
	return ul.New(items...).Class("diagnostic-list")
}

// Event renders a single diagnostic event entry in the live feed.
func Event(s string) node.Node {
	return li.New().Class("diagnostic-item").Text(s)
}

// Package toc provides table-of-contents elements for section overview
// pages - a styled list with linked entries and feature summaries.
package toc

import (
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"
)

// List creates a styled table-of-contents container.
func List(items ...node.Node) node.Node {
	return ul.New(items...).Class("toc")
}

// Item renders a single table-of-contents entry with a link and a
// feature summary line.
func Item(link node.Node, features string) node.Node {
	return li.New(
		link,
		span.Text(features).Class("toc-features"),
	).Class("toc-item")
}

// Link creates a plain anchor for table-of-contents entries - used
// in the HTTP section where there is no persistent session.
func Link(path, title string) node.Node {
	return a.New().Href(path).Class("toc-title").Text(title)
}

// NavLink creates a client-side navigation link for table-of-contents
// entries - used in the WebSocket/SSE sections where bind.Link()
// navigates without a full page reload.
func NavLink(path, title string) node.Node {
	return bind.Apply(a.New().Href(path).Class("toc-title").Text(title), bind.Link())
}

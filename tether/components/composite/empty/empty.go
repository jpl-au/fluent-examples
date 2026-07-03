// Package empty provides a styled empty-state page for 404 and
// placeholder content.
package empty

import (
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h2"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"
)

// State renders a centred empty-state block with a title, message,
// and an optional link back to a parent page.
func State(title, message string, link node.Node) node.Node {
	return div.New(
		h2.Text(title).Class("empty-title"),
		p.Text(message).Class("empty-text"),
		node.When(link != nil, link),
	).Class("empty-state")
}

// Link creates a plain anchor for the empty-state back link.
func Link(path, text string) node.Node {
	return a.Text(text).Href(path).Class("empty-link")
}

// NavLink creates a client-side navigation link for the empty-state
// back link.
func NavLink(path, text string) node.Node {
	return bind.Apply(a.Text(text).Href(path).Class("empty-link"), bind.Link())
}

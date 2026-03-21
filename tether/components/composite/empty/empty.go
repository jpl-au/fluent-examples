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
	children := []node.Node{
		h2.New().Class("empty-title").Text(title),
		p.New().Class("empty-text").Text(message),
	}
	if link != nil {
		children = append(children, link)
	}
	return div.New(children...).Class("empty-state")
}

// Link creates a plain anchor for the empty-state back link.
func Link(path, text string) node.Node {
	return a.New().Href(path).Class("empty-link").Text(text)
}

// NavLink creates a client-side navigation link for the empty-state
// back link.
func NavLink(path, text string) node.Node {
	return bind.Apply(a.New().Href(path).Class("empty-link").Text(text), bind.Link())
}

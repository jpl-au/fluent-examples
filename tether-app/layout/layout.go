// Package layout provides the page shell for the kanban board
// application: a header with the app title, user name, online count,
// and an action area, wrapping the main content region.
package layout

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h1"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"
)

// Shell wraps content in the app chrome: a header bar and scrollable
// content area. The online count badge updates reactively via a
// signal binding.
func Shell(name string, onlineCount int, actions node.Node, content node.Node) node.Node {
	return div.New(
		header(name, onlineCount, actions),
		div.New(content).Class("content").Dynamic("_"),
	).Class("shell")
}

// header builds the top bar with app title, user name, online badge,
// and action nodes.
func header(name string, onlineCount int, actions node.Node) node.Node {
	badge := bind.Apply(
		span.New().Class("badge badge-green").Text(fmt.Sprintf("%d online", onlineCount)),
		bind.BindText("online_count"),
	)

	left := div.New(
		h1.New().Text("Kanban Board"),
		span.New().Class("header-user").Text(name),
		div.New(badge).Class("header-meta"),
	).Class("header-left")

	return div.New(left, div.New(actions).Class("header-actions")).Class("header").Dynamic("header")
}

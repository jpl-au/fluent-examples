package page

import (
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/empty"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// NotFoundRender builds the 404 page for the Service Worker section,
// shown when the router has no match for the current URL path.
func NotFoundRender(_ state.State) node.Node {
	return cpage.New(
		empty.State("Page not found", "The page you're looking for doesn't exist in this section.",
			empty.NavLink("/sw/", "Back to overview"),
		),
	)
}

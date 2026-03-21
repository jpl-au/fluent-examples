package http

import (
	"github.com/jpl-au/tether/router"

	"github.com/jpl-au/fluent-examples/tether/site/http/page"
	"github.com/jpl-au/fluent-examples/tether/site/http/state"
)

// newRouter registers all URL paths for the stateless HTTP section.
func newRouter() *router.Router[state.State] {
	r := router.New[state.State](func(s state.State) string { return s.Page })

	r.Route("/", router.Page[state.State]{
		Render: page.OverviewRender,
	})
	r.NotFound(router.Page[state.State]{
		Render: page.NotFoundRender,
	})

	return r
}

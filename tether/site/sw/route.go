package sw

import (
	"github.com/jpl-au/tether/router"

	"github.com/jpl-au/fluent-examples/tether/site/sw/page"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// newRouter registers all URL paths for the Service Worker section.
func newRouter() *router.Router[state.State] {
	r := router.New[state.State](func(s state.State) string { return s.Page })

	r.Route("/sw/", router.Page[state.State]{Render: page.OverviewRender})
	r.Route("/sw/push", router.Page[state.State]{Render: page.PushRender, Handle: page.PushHandle})
	r.Route("/sw/caching", router.Page[state.State]{Render: page.CachingRender})
	r.Route("/sw/lifecycle", router.Page[state.State]{Render: page.LifecycleRender, Handle: page.LifecycleHandle})
	r.NotFound(router.Page[state.State]{Render: page.NotFoundRender})

	return r
}

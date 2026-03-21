// Package menu provides a three-dot dropdown menu for contextual
// actions. Navigation items use HTMX for partial swaps; destructive
// actions use hx-post on a form.
package menu

import (
	htmx "github.com/jpl-au/fluent-htmx"
	"github.com/jpl-au/fluent-htmx/swap"
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/node"
)

// New creates a three-dot dropdown menu with inline toggle script.
func New(items ...node.Node) node.Node {
	return div.New(
		div.Static("⋮").Class("menu-trigger"),
		div.New(items...).Class("menu-dropdown"),
		toggleScript(),
	).Class("menu")
}

// Link creates a navigation menu item that uses HTMX to swap the
// content area without a full page reload.
func Link(label, href string) node.Node {
	link := a.Text(label).Href(href).Class("menu-item")
	htmx.New(link).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)
	return link
}

// DangerLink creates a destructive navigation menu item with HTMX.
func DangerLink(label, href string) node.Node {
	link := a.Text(label).Href(href).Class("menu-item menu-item-danger")
	htmx.New(link).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)
	return link
}

// FormAction creates a menu item that submits a POST via HTMX  -
// used for destructive actions like Delete.
func FormAction(label, action string) node.Node {
	f := form.Post(action)
	htmx.New(f).HxPost(action).HxTarget("#content").HxSwap(swap.InnerHTML)
	return div.New(
		f.Add(button.Submit(label).Class("menu-item menu-item-danger")),
	).Class("menu-form")
}

// toggleScript returns the inline script for menu toggle behaviour.
func toggleScript() node.Node {
	return script.Static(`
(function() {
  var s = document.currentScript;
  var menu = s.closest('.menu');
  if (!menu) return;
  var trigger = menu.querySelector('.menu-trigger');
  var dropdown = menu.querySelector('.menu-dropdown');

  trigger.addEventListener('click', function(e) {
    e.stopPropagation();
    dropdown.classList.toggle('menu-open');
  });

  document.addEventListener('click', function() {
    dropdown.classList.remove('menu-open');
  });

  dropdown.addEventListener('click', function(e) {
    e.stopPropagation();
  });
})();
`)
}

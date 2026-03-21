// Package menu provides a three-dot dropdown menu for contextual
// actions. The toggle is handled by a small inline script - no
// external JS dependencies. Click the trigger to open, click
// anywhere else (or the trigger again) to close.
package menu

import (
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/node"
)

// New creates a three-dot dropdown menu. Items are the menu entries
// - typically built with Link and Action. The script for toggling
// is embedded inline so the menu works without any external JS.
func New(items ...node.Node) node.Node {
	return div.New(
		div.Static("⋮").Class("menu-trigger"),
		div.New(items...).Class("menu-dropdown"),
		toggleScript(),
	).Class("menu")
}

// Link creates a navigation menu item that links to a URL.
func Link(label, href string) node.Node {
	return a.Text(label).Href(href).Class("menu-item")
}

// DangerLink creates a destructive navigation menu item.
func DangerLink(label, href string) node.Node {
	return a.Text(label).Href(href).Class("menu-item menu-item-danger")
}

// FormAction creates a menu item that submits a POST form - used for
// destructive actions like Delete that must not be GET requests.
func FormAction(label, action string) node.Node {
	return form.Post(action,
		button.Submit(label).Class("menu-item menu-item-danger"),
	).Class("menu-form")
}

// toggleScript returns the inline script that wires up the menu
// toggle behaviour. Clicking the trigger toggles the dropdown;
// clicking outside closes it. Each menu instance is self-contained
// - the script walks up from itself to find its parent .menu.
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

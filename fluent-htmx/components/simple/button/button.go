// Package button provides styled button and link elements. Plain
// variants navigate via standard GET; Hx variants use HTMX to swap
// the content area without a full page reload.
package button

import (
	htmx "github.com/jpl-au/fluent-htmx"
	"github.com/jpl-au/fluent-htmx/swap"
	"github.com/jpl-au/fluent/html5/a"
	el "github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/node"
)

// Primary creates a high-contrast link styled as a button. Uses HTMX
// to swap the #content area and push the URL to browser history.
func Primary(label, href string) node.Node {
	link := a.Text(label).Href(href).Class("btn btn-primary")
	htmx.New(link).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)
	return link
}

// Secondary creates a medium-contrast HTMX navigation button.
func Secondary(label, href string) node.Node {
	link := a.Text(label).Href(href).Class("btn btn-secondary")
	htmx.New(link).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)
	return link
}

// Danger creates a destructive-action HTMX navigation button.
func Danger(label, href string) node.Node {
	link := a.Text(label).Href(href).Class("btn btn-danger")
	htmx.New(link).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)
	return link
}

// Link creates a plain styled HTMX navigation link.
func Link(label, href string) node.Node {
	link := a.Text(label).Href(href).Class("btn btn-link")
	htmx.New(link).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)
	return link
}

// Back creates a secondary HTMX navigation button with a left arrow.
func Back(href string) node.Node {
	link := a.Static("← Back").Href(href).Class("btn btn-secondary")
	htmx.New(link).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)
	return link
}

// Submit creates a form submit button. No HTMX - forms use hx-post
// on the form element itself.
func Submit(label string) node.Node {
	return el.Submit(label).Class("btn btn-primary")
}

// DangerSubmit creates a destructive form submit button.
func DangerSubmit(label string) node.Node {
	return el.Submit(label).Class("btn btn-danger btn-sm")
}

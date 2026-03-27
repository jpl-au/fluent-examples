// Package layout provides the shell for the Tether feature
// explorer: sidebar navigation with feature groups, header, and
// content area.
package layout

import (
	"fmt"

	"github.com/jpl-au/fluent-examples/tether/components/simple/badge"
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h1"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/nav"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"
)

// Section identifies the handler type for the current page.
type Section int

const (
	// SectionHTTP is for stateless tether.Stateless handlers.
	SectionHTTP Section = iota
	// SectionLive is for standalone tether.Handler instances.
	SectionLive
	// SectionSW is for the service worker handler, which uses
	// internal client-side navigation between its sub-pages.
	SectionSW
)

// Live reports whether the section uses a persistent connection.
func (s Section) Live() bool {
	return s != SectionHTTP
}

// sidebarGroup holds a top-level sidebar heading and its items.
type sidebarGroup struct {
	heading   string
	subgroups []sectionGroup
}

// sectionGroup holds navigation metadata for one section.
type sectionGroup struct {
	section Section
	heading string // rendered as a divider label; empty to omit
	items   []navItem
}

// navItem describes a single sidebar link.
type navItem struct {
	path  string
	label string
}

// sidebar defines the navigation groups, organised by feature.
var sidebar = []sidebarGroup{
	{"Stateless (tether.Stateless)", []sectionGroup{
		{SectionHTTP, "", []navItem{
			{"/", "Overview"},
			{"/events", "Events & Forms"},
			{"/rendering", "State & Rendering"},
			{"/errors", "Error Boundaries"},
			{"/navigation", "Navigation"},
			{"/middleware", "Middleware"},
			{"/morph", "Full-Page Morph"},
			{"/clipboard", "Clipboard"},
			{"/selection", "Multi-Select"},
			{"/touch", "Touch Gestures"},
		}},
	}},
	{"Signals & Directives", []sectionGroup{
		{SectionLive, "", []navItem{
			{"/signals/ws/", "WebSocket"},
			{"/signals/sse/", "SSE"},
		}},
	}},
	{"Live Updates", []sectionGroup{
		{SectionLive, "", []navItem{
			{"/live/ws/", "WebSocket"},
			{"/live/sse/", "SSE"},
		}},
	}},
	{"Features", []sectionGroup{
		{SectionLive, "", []navItem{
			{"/notifications/", "Notifications"},
			{"/uploads/", "File Uploads"},
			{"/uploads/filtered/", "Filtered Uploads"},
			{"/broadcasting/", "Broadcasting"},
			{"/components/", "Components"},
			{"/chat/", "Chat Room"},
			{"/realtime/", "Real-time Dashboard"},
			{"/configuration/", "Configuration"},
			{"/valuestore/", "Value Store"},
			{"/groups/", "Groups"},
			{"/freeze/", "Freeze & Restore"},
			{"/hotkey/", "Hotkeys"},
			{"/dragdrop/", "Drag and Drop"},
			{"/scroll/", "Scroll"},
		}},
	}},
	{"Performance", []sectionGroup{
		{SectionLive, "", []navItem{
			{"/memo/", "Memoisation"},
			{"/memo/realtime/", "Memoised Dashboard"},
		}},
	}},
	{"Observability", []sectionGroup{
		{SectionLive, "", []navItem{
			{"/diagnostics/", "Diagnostics"},
		}},
	}},
	{"Service Worker", []sectionGroup{
		{SectionSW, "", []navItem{
			{"/sw/", "Overview"},
			{"/sw/push", "Push Notifications"},
			{"/sw/caching", "Caching & Offline"},
		}},
	}},
}

// Shell wraps page content in the app chrome: sidebar and header.
// The section parameter determines whether the online count badge
// appears. The content wrapper uses Dynamic("_") - a pass-through
// marker that lets the diff engine look through to the page-level
// keys inside so each region is patched independently.
func Shell(section Section, currentPage string, onlineCount int, content node.Node) node.Node {
	return div.New(
		sidebarNav(section, currentPage),
		div.New(
			header(section, currentPage, onlineCount),
			div.New(content).Class("content").Dynamic("_"),
		).Class("main"),
	).Class("shell")
}

// sidebarNav builds the full sidebar from the static sidebar data.
func sidebarNav(section Section, currentPage string) node.Node {
	var items []node.Node
	for _, sg := range sidebar {
		items = append(items, li.New(
			span.Text(sg.heading).Class("nav-group"),
		))
		for _, sub := range sg.subgroups {
			if sub.heading != "" {
				items = append(items, li.New(
					span.Text(sub.heading).Class("nav-divider"),
				))
			}
			for _, item := range sub.items {
				items = append(items, navLink(section, sub.section, item, currentPage))
			}
		}
	}

	return nav.New(
		div.New(
			h1.Static("Tether"),
			span.Static("Feature Explorer"),
		).Class("sidebar-brand"),
		ul.New(items...).Class("sidebar-nav"),
		sidebarScrollScript(),
	).Class("sidebar").Dynamic("sidebar")
}

// sameHandler reports whether two sections are served by the same
// tether.Handler. Only the service worker handler has internal
// navigation between sub-pages; every other live handler is
// standalone, so cross-feature links always do full page loads.
func sameHandler(a, b Section) bool {
	return a == SectionSW && b == SectionSW
}

// navLink renders a sidebar link. Links within the same handler
// use bind.Link for client-side navigation. All other links use
// plain <a> tags for full page loads.
func navLink(currentSection, linkSection Section, item navItem, currentPage string) node.Node {
	cls := "nav-link"
	if item.path == currentPage {
		cls += " active"
	}

	link := a.New().Href(item.path).Class(cls).Add(
		span.Text(item.label),
	)

	if sameHandler(currentSection, linkSection) && currentSection.Live() {
		return li.New(bind.Apply(link, bind.Link()))
	}
	return li.New(link)
}

// header builds the page header with an optional online count badge.
func header(section Section, currentPage string, onlineCount int) node.Node {
	title := pageTitle(currentPage)
	headerNode := h1.Text(title)

	if section.Live() {
		return div.New(
			headerNode,
			div.New(
				bind.Apply(
					badge.Green(fmt.Sprintf("%d online", onlineCount)),
					bind.BindText("online_count"),
				),
			).Class("header-actions"),
		).Class("header").Dynamic("header")
	}

	return div.New(headerNode).Class("header").Dynamic("header")
}

// pageTitle looks up the human-readable label for a URL path by
// walking the sidebar navigation tree.
func pageTitle(page string) string {
	for _, sg := range sidebar {
		for _, sub := range sg.subgroups {
			for _, item := range sub.items {
				if item.path == page {
					return item.label
				}
			}
		}
	}
	return "Not Found"
}

// sidebarScrollScript returns an inline script that persists the
// sidebar scroll position across full page loads using sessionStorage.
// Without this, navigating between handlers (which triggers a full
// page load) resets the sidebar to the top, losing the user's place.
func sidebarScrollScript() node.Node {
	return script.Static(`(function(){
  var sb = document.querySelector('.sidebar');
  if (!sb) return;
  var key = '_sidebar_scroll';
  var saved = sessionStorage.getItem(key);
  if (saved) sb.scrollTop = parseInt(saved, 10);
  sb.addEventListener('scroll', function(){
    sessionStorage.setItem(key, sb.scrollTop);
  });
})();`)
}

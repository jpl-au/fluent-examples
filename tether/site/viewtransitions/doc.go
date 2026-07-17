// Package viewtransitions demonstrates native View Transitions. The
// handler opts in with tether.Client{ViewTransitions: true}, so every
// server-driven DOM update is wrapped in document.startViewTransition
// and cross-fades instead of snapping. A shared element carries a CSS
// view-transition-name so the browser animates it between positions.
// Stateless HTTP - the effect applies to ordinary morph updates.
package viewtransitions

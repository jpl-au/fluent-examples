package handler

import (
	security "github.com/jpl-au/fluent-security"
	"github.com/microcosm-cc/bluemonday"
)

// descriptionPolicy permits the bluemonday UGC baseline plus `class`
// attributes on <code>/<pre>. Rich text pasted from editors that
// syntax-highlight snippets (e.g. highlight.js output) keeps its
// styling hooks intact; scripts, event handlers, and javascript:
// URIs are still dropped by the UGC baseline.
var descriptionPolicy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("code", "pre")
	return p
}()

// descCleaner is the hoisted sanitiser used for every card
// description. Components that need to render an untrusted
// description receive this value rather than constructing their own,
// so the policy is defined once and every call site sees the same
// rules.
var descCleaner = security.FromPolicy(descriptionPolicy)

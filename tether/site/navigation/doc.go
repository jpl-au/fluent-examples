// Package navigation demonstrates client-side navigation with
// Tether: bind.Link for pushState navigation, OnNavigate for query
// parameter extraction (Params.Get, Params.IntDefault, Params.BoolDefault,
// Params.Float64Default, Params.Strings, Params.Ints, Params.Float64s),
// Session.ReplaceURL for history-silent URL updates, and
// Session.Navigate for server-driven navigation.
//
// This example uses tether.Stateless with a small internal router so the
// bind.Link demos have pages to navigate between.
package navigation

// Package htmlwire demonstrates the HTML wire format for stateless
// pages: POST event responses are plain HTML - the morph fragments as
// the response body, side effects in a JSON island - instead of the
// default JSON envelope. It also shows StatelessConfig.CacheControl,
// which makes the (session-token-free) initial GET cacheable.
//
// Exercises: StatelessConfig.WireFormat (wire.HTML),
// StatelessConfig.CacheControl, Session.Morph targeted fragments, and
// the Tether-Morph response header.
package htmlwire

// Package handler provides HTTP request handlers for the contact
// manager example application. It covers three communication styles:
//
//   - Contact CRUD handlers that use htmx partial responses. GET
//     handlers check htmx.HxRequestType(r) to decide between a full
//     page and a content-only partial; a history restore asks for the
//     full page. POST handlers use htmx.Handle(r, func(){}) to return
//     updated partials instead of redirecting. Partial responses carry
//     a <title> and <hx-partial> blocks so the tab title and the
//     header update alongside the content.
//
//   - A WebSocket handler (hx-ws:connect, in the htmax.js bundle) that
//     upgrades the connection and pushes rendered HTML log entries as
//     plain HTML, which htmx swaps in with the connecting element's
//     own target and swap style. The first message is a JSON frame
//     that names its own target. Lines the page sends with hx-ws:send
//     come back on the same connection as entries.
//
//   - A Server-Sent Events handler (hx-sse:connect, in the htmax.js
//     bundle) that streams rendered HTML log entries as unnamed
//     events. htmx swaps each one into the connecting element with its
//     own swap style.
package handler

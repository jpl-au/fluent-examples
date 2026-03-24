// Package handler provides HTTP request handlers for the contact
// manager example application. It covers three communication styles:
//
//   - Contact CRUD handlers that use htmx partial responses. GET
//     handlers check htmx.HxRequest(r) to decide between a full page
//     and a content-only partial. POST handlers use htmx.Handle(r,
//     func(){}) to return updated partials instead of redirecting.
//     Partial responses include an out-of-band header swap so the
//     page title and header actions update alongside the content.
//
//   - A WebSocket handler (htmx-ext-ws) that upgrades the connection
//     and pushes rendered HTML log entries as OOB swap fragments. The
//     htmx WebSocket extension on the client receives these fragments
//     and swaps them into the page.
//
//   - A Server-Sent Events handler (htmx-ext-sse) that streams
//     rendered HTML log entries over an EventSource connection. The
//     htmx SSE extension on the client receives named events and
//     swaps the HTML payloads into matching elements.
package handler

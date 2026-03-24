// Package handler provides HTTP request handlers for the contact
// manager application. It covers three concerns:
//
//   - Contact and note CRUD using the PRG (Post-Redirect-Get)
//     pattern: GET handlers render pages, POST handlers mutate
//     state and redirect.
//   - A WebSocket handler that pushes live-generated events to
//     connected clients.
//   - An SSE (Server-Sent Events) handler that streams
//     live-generated events over a single long-lived HTTP response.
package handler

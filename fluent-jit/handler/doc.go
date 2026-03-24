// Package handler provides request handlers for the contact manager
// application. It covers three transport mechanisms:
//
//   - HTTP handlers following the PRG (Post-Redirect-Get) pattern for
//     standard page navigation and form submissions.
//   - A WebSocket handler that upgrades connections and pushes live log
//     entries over a full-duplex channel.
//   - An SSE handler that streams log entries as Server-Sent Events
//     over a unidirectional connection.
//
// Across these handlers, the package demonstrates all three fluent-jit
// strategies (Compile, Flatten, and Tune) using both the Global API
// (string-keyed registry) and the Instance API (fine-grained control).
package handler

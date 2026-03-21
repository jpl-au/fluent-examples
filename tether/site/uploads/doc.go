// Package uploads demonstrates file uploads via bind.Upload with
// real-time feedback: the upload handler uses the live session to
// push a confirmation toast and append the file to the rendered list
// immediately via sess.Update. Transport-agnostic - the upload itself
// is always HTTP POST; only the feedback channel requires WebSocket.
package uploads

// Package freeze demonstrates FreezeWithConnect: when the client
// disconnects, the server persists session state to the SessionStore,
// releases memory, and exits the command loop. On reconnect, state is
// restored from the store - the counter picks up where it left off.
package freeze

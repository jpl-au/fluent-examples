// Package diagnostics demonstrates Handler.Diagnostics - the runtime
// event bus that emits framework-level events such as transport errors,
// buffer overflows, handler panics, and upload rejections. Events are
// pushed into session state via WatchBus and rendered as a live feed,
// showing developers how to build observability into their application.
package diagnostics

// Package patch demonstrates targeted updates with sess.Patch.
// Patch works with either engine (Differ or Memoiser) and does not
// require Memoise: true. A background timer increments a single counter
// in a list of 20, re-rendering only the affected row instead of
// the full page. Each update takes ~5-10µs vs ~3-12ms for a full
// render cycle.
package patch

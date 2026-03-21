// Package morph demonstrates the full-page morph fallback that occurs
// when no Dynamic keys are present. The differ finds no keyed elements
// to patch, so the framework sends the entire rendered HTML and the
// client-side idiomorph library diffs the whole DOM. This is slower
// than targeted patches but still produces correct updates. Stateless
// HTTP - no persistent connection needed.
package morph

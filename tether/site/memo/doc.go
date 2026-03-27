// Package memo demonstrates subtree memoisation with tether.Versioned
// and node.Memo. An expensive table region is skipped when its data
// hasn't changed, while a cheap counter updates on every click.
package memo

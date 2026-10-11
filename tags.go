//go:build linux && !android

package antui

// A tags is a tagstruct being built: values, each behind the byte that says
// what it is.
type tags struct{ b []byte }

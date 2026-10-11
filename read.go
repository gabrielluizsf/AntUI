//go:build linux && !android

package antui

// A read is a tagstruct being taken apart.
type read struct {
	b  []byte
	at int
}

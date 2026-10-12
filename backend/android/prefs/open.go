//go:build android

package prefs

// Open names a store. Nothing happens until something is read or written, so
// this cannot fail.
func Open(name string) *Prefs { return &Prefs{name: name} }

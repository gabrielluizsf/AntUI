//go:build android

package ndk

// Has reports whether every modifier in want is held.
func (m Meta) Has(want Meta) bool { return m&want == want }

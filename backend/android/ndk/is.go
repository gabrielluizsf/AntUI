//go:build android

package ndk

// Is reports whether the source includes want. A touchscreen and a stylus
// both carry the touchscreen bits, so this is an "includes", not an "equals".
func (s Source) Is(want Source) bool { return s&want == want }

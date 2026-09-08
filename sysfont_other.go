//go:build !windows && !darwin && !linux && !freebsd && !netbsd && !openbsd && !dragonfly

package antui

// defaultUIPoints is a guess, on a platform this library has not been taught
// about.
const defaultUIPoints = 11

// systemFontAsked has nothing to ask on a platform this library has not
// been taught about.
func systemFontAsked() []string { return nil }

// systemFontPaths is empty here, so [SystemFace] answers ErrNoSystemFont and
// the program carries on with the built-in — which is the right outcome and
// not a failure.
func systemFontPaths() []string { return nil }
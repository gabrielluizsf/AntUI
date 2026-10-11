//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// Predefined atoms the protocol fixes, so they need no InternAtom round trip.
const (
	atomString   = 31
	atomAtom     = 4
	atomCardinal = 6
	atomWMName   = 39
	atomWMHints  = 35
	atomWMClass  = 67
)

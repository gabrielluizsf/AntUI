//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// keysymToKey maps a symbol to the physical key it stands for. Printable
// keys report their uppercase ASCII code, so the caller can compare against
// KeyA regardless of the layout's case.
func keysymToKey(sym uint32) backend.Key {
	switch {
	case sym >= 'a' && sym <= 'z':
		return backend.Key(sym - 32)
	case sym >= 'A' && sym <= 'Z', sym >= '0' && sym <= '9':
		return backend.Key(sym)
	}
	if key := keysymNav(sym); key != backend.KeyUnknown {
		return key
	}
	if key := keysymMod(sym); key != backend.KeyUnknown {
		return key
	}
	return keysymRange(sym)
}

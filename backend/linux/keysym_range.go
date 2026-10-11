//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// keysymRange maps the two families whose symbols run consecutively, the
// function keys and the numeric keypad, and treats everything below 0x80 that
// was not named as ASCII punctuation.
func keysymRange(sym uint32) backend.Key {
	switch {
	case sym >= symF1 && sym <= symF12:
		return backend.KeyF1 + backend.Key(sym-symF1)
	case sym >= symKP0 && sym <= symKP9:
		return backend.Key0 + backend.Key(sym-symKP0)
	case sym < 0x80:
		return backend.Key(sym) // ASCII punctuation
	}
	return backend.KeyUnknown
}

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// keysymMod maps the modifier keys.
func keysymMod(sym uint32) backend.Key {
	switch sym {
	case symShiftL:
		return backend.KeyLeftShift
	case symShiftR:
		return backend.KeyRightShift
	case symCtrlL:
		return backend.KeyLeftControl
	case symCtrlR:
		return backend.KeyRightControl
	case symCaps:
		return backend.KeyCapsLock
	case symAltL:
		return backend.KeyLeftAlt
	case symAltR, symISOLevel3:
		return backend.KeyRightAlt
	case symSuperL:
		return backend.KeyLeftSuper
	case symSuperR:
		return backend.KeyRightSuper
	}
	return backend.KeyUnknown
}

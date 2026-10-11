//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// winKeyMods maps the modifier keys, where the extended bit is what tells the
// two sides apart for the keys Windows reports under one code.
func winKeyMods(vk uintptr, extended bool) (backend.Key, bool) {
	switch vk {
	case vkControl:
		if extended {
			return backend.KeyRightControl, true
		}
		return backend.KeyLeftControl, true
	case vkMenu:
		if extended {
			return backend.KeyRightAlt, true
		}
		return backend.KeyLeftAlt, true
	case vkLCtrl:
		return backend.KeyLeftControl, true
	case vkRCtrl:
		return backend.KeyRightControl, true
	case vkLMenu:
		return backend.KeyLeftAlt, true
	case vkRMenu:
		return backend.KeyRightAlt, true
	case vkLWin:
		return backend.KeyLeftSuper, true
	case vkRWin:
		return backend.KeyRightSuper, true
	}
	return 0, false
}

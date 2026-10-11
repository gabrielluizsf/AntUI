//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// winKeyNav maps the navigation, editing and lock keys.
func winKeyNav(vk uintptr, extended bool) (backend.Key, bool) {
	switch vk {
	case vkSpace:
		return backend.KeySpace, true
	case vkEscape:
		return backend.KeyEscape, true
	case vkReturn:
		return backend.KeyEnter, true
	case vkTab:
		return backend.KeyTab, true
	case vkBack:
		return backend.KeyBackspace, true
	case vkInsert:
		return backend.KeyInsert, true
	case vkDelete:
		return backend.KeyDelete, true
	case vkRight:
		return backend.KeyRight, true
	case vkLeft:
		return backend.KeyLeft, true
	case vkDown:
		return backend.KeyDown, true
	case vkUp:
		return backend.KeyUp, true
	case vkPrior:
		return backend.KeyPageUp, true
	case vkNext:
		return backend.KeyPageDown, true
	case vkHome:
		return backend.KeyHome, true
	case vkEnd:
		return backend.KeyEnd, true
	case vkCapital:
		return backend.KeyCapsLock, true
	case vkShift, vkLShift:
		return backend.KeyLeftShift, true
	case vkRShift:
		return backend.KeyRightShift, true
	}
	return 0, false
}

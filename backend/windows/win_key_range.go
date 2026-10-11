//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// winKeyRange maps the two families whose codes run consecutively: the
// function keys and the numeric keypad.
func winKeyRange(vk uintptr) backend.Key {
	switch {
	case vk >= vkF1 && vk <= vkF12:
		return backend.KeyF1 + backend.Key(vk-vkF1)
	case vk >= vkNumpad0 && vk <= vkNumpad9:
		return backend.Key0 + backend.Key(vk-vkNumpad0)
	}
	return backend.KeyUnknown
}

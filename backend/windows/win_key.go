//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// winKey maps a virtual key code to the key it stands for. The extended bit
// in lparam is what separates the left and right control and alt keys.
func winKey(vk, lparam uintptr) backend.Key {
	extended := lparam&(1<<24) != 0

	if vk >= 'A' && vk <= 'Z' || vk >= '0' && vk <= '9' {
		return backend.Key(vk)
	}
	if key, ok := winKeyNav(vk, extended); ok {
		return key
	}
	if key, ok := winKeyMods(vk, extended); ok {
		return key
	}
	if key, ok := winKeyPunct(vk); ok {
		return key
	}
	return winKeyRange(vk)
}

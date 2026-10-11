//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// winKeyPunct maps the OEM punctuation keys, whose codes vary by layout but
// whose names here do not.
func winKeyPunct(vk uintptr) (backend.Key, bool) {
	switch vk {
	case vkOEMComma:
		return backend.KeyComma, true
	case vkOEMPeriod:
		return backend.KeyPeriod, true
	case vkOEMMinus:
		return backend.KeyMinus, true
	case vkOEMPlus:
		return backend.KeyEqual, true
	}
	return 0, false
}
